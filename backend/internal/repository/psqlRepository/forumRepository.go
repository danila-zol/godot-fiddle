package psqlRepository

import (
	"context"
	"errors"
	"fmt"
	"gamehangar/internal/domain/models"
	"strings"

	"github.com/jackc/pgx/v5"
)

type PsqlForumRepository struct {
	databaseClient psqlDatabaseClient
	enforcer       Enforcer
	conflictErr    error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlForumRepository(dbClient psqlDatabaseClient, e Enforcer) *PsqlForumRepository {
	return &PsqlForumRepository{
		databaseClient: dbClient,
		enforcer:       e,
		conflictErr:    errors.New("Record conflict!"),
	}
}

func (r *PsqlForumRepository) NotFoundErr() error { return r.databaseClient.ErrNoRows() }

// Returns "Record conflict!" to specify conflicting record versions on update
func (r *PsqlForumRepository) ConflictErr() error { return r.conflictErr }

func (r *PsqlForumRepository) CreateTopic(topic models.Topic) (*models.Topic, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO forum.topics
		(name) 
		VALUES
		($1)
		RETURNING
		(id, name, version)`,
		topic.Name,
	).Scan(&topic)

	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *PsqlForumRepository) FindTopicByID(id int) (*models.Topic, error) {
	var topic models.Topic
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT (id, name, version) FROM forum.topics WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *PsqlForumRepository) FindTopics() (*[]models.Topic, error) {
	var topics []models.Topic

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT (id, name, version) FROM forum.topics`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var topic models.Topic
		err = rows.Scan(&topic)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(topics) == 0 {
		return nil, r.NotFoundErr()
	}
	return &topics, nil
}

func (r *PsqlForumRepository) UpdateTopic(id int, topic models.Topic) (*models.Topic, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(
		context.Background(),
		`SELECT (id) FROM forum.topics WHERE id=$1 AND version=$2`,
		id, *topic.Version,
	).Scan(&id)
	if err != nil {
		return nil, r.ConflictErr()
	}

	err = conn.QueryRow(context.Background(),
		`UPDATE forum.topics SET 
		name=COALESCE($1, name)
		WHERE id = $2
		RETURNING
		(id, name, version)`,
		topic.Name, id,
	).Scan(&topic)
	if err != nil {
		return nil, err
	}
	return &topic, err
}

func (r *PsqlForumRepository) DeleteTopic(id int) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	err = r.DeleteThreadsOfTopic(id)
	if err != nil {
		return err
	}

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.topics WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	return nil
}

func (r *PsqlForumRepository) CreateThread(thread models.Thread) (*models.Thread, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO forum.threads
			(title, user_id, topic_id, tags)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			(id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes, rating, views)`,
		thread.Title, thread.UserID, thread.TopicID, thread.Tags,
	).Scan(&thread)
	if err != nil {
		return nil, err
	}

	_, err = r.enforcer.AddPermissions(thread.UserID.String(), fmt.Sprintf("threads/%v", *thread.ID), "PATCH")
	if err != nil {
		return nil, err
	}
	_, err = r.enforcer.AddPermissions(thread.UserID.String(), fmt.Sprintf("threads/%v", *thread.ID), "DELETE")
	if err != nil {
		return nil, err
	}

	return &thread, err
}

func (r *PsqlForumRepository) FindThreadByID(id int) (*models.Thread, error) {
	var thread models.Thread
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE forum.threads SET 
		views=views+1
		WHERE id = $1
		RETURNING
		(id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes, rating, views)`,
		id,
	).Scan(&thread)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r *PsqlForumRepository) FindThreads(keywords []string, limit uint64, order string) (*[]models.Thread, error) {
	var threads []models.Thread

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var rows pgx.Rows
	if len(keywords) != 0 {
		query := `SELECT (id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes, rating, views) 
				FROM
				((SELECT id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes, rating, views
				FROM forum.threads
				WHERE thread_ts @@ to_tsquery_multilang($1))
			UNION
				(SELECT id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes, rating, views 
				FROM forum.threads
				WHERE tags && ($2) COLLATE case_insensitive))`

		switch order {
		case "newest-updated":
			query = query + ` ORDER BY updated_at DESC`
		case "highest-rated":
			query = query + ` ORDER BY rating DESC`
		case "most-views":
			query = query + ` ORDER BY views DESC`
		default:
			query = query + ` ORDER BY updated_at DESC`
		}
		if limit != 0 {
			query = query + fmt.Sprintf(` LIMIT %v`, limit)
		}
		rows, err = conn.Query(context.Background(),
			query, strings.Join(keywords, " | "), keywords,
		)
		if err != nil {
			return nil, err
		}
	} else {
		query := `SELECT (id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes, rating, views)
		FROM forum.threads`

		switch order {
		case "newest-updated":
			query = query + ` ORDER BY updated_at DESC`
		case "highest-rated":
			query = query + ` ORDER BY rating DESC`
		case "most-views":
			query = query + ` ORDER BY views DESC`
		default:
			query = query + ` ORDER BY updated_at DESC`
		}
		if limit != 0 {
			query = query + fmt.Sprintf(` LIMIT %v`, limit)
		}
		rows, err = conn.Query(context.Background(), query)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()
	for rows.Next() {
		var thread models.Thread
		err = rows.Scan(&thread)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(threads) == 0 {
		return nil, r.NotFoundErr()
	}
	return &threads, nil
}

func (r *PsqlForumRepository) UpdateThread(id int, thread models.Thread) (*models.Thread, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE forum.threads SET 
			title=COALESCE($1, title), user_id=COALESCE($2, user_id), topic_id=COALESCE($3, topic_id),
		tags=COALESCE($4, tags), upvotes=COALESCE($5, upvotes), downvotes=COALESCE($6, downvotes),
			updated_at=NOW()
			WHERE id = $7
		RETURNING
			(id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes, rating, views)`,
		thread.Title, thread.UserID, thread.TopicID, thread.Tags,
		thread.Upvotes, thread.Downvotes, id,
	).Scan(&thread)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r *PsqlForumRepository) DeleteThread(id int) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.threads WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	_, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("threads/%v", id), "PATCH")
	if err != nil {
		return err
	}
	_, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("threads/%v", id), "DELETE")
	if err != nil {
		return err
	}
	return nil
}

func (r *PsqlForumRepository) DeleteThreadsOfTopic(topicID int) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT id FROM forum.threads WHERE topic_id = $1`, topicID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var thread int
		err = rows.Scan(&thread)
		if err != nil {
			return err
		}
		err = r.DeleteMessagesOfThread(thread)
		if err != nil {
			return err
		}
		_, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("threads/%v", thread), "PATCH")
		if err != nil {
			return err
		}
		_, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("threads/%v", thread), "DELETE")
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *PsqlForumRepository) CreateMessage(message models.Message) (*models.Message, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO forum.messages
		(thread_id, user_id, title, body, tags) 
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING
		(id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes, rating, views)`,
		message.ThreadID, message.UserID, message.Title, message.Body, message.Tags,
	).Scan(&message)
	if err != nil {
		return nil, err
	}

	// TODO: When deleting a Forum it leaves the Thread and Message permissions intact
	_, err = r.enforcer.AddPermissions(message.UserID.String(), fmt.Sprintf("messages/%v", *message.ID), "PATCH")
	if err != nil {
		return nil, err
	}
	_, err = r.enforcer.AddPermissions(message.UserID.String(), fmt.Sprintf("messages/%v", *message.ID), "DELETE")
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *PsqlForumRepository) FindMessageByID(id int) (*models.Message, error) {
	var message models.Message
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE forum.messages SET 
		views=views+1
		WHERE id = $1
		RETURNING
		(id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes, rating, views)`,
		id,
	).Scan(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *PsqlForumRepository) FindMessages(keywords []string, limit uint64, order string) (*[]models.Message, error) {
	var messages []models.Message

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var rows pgx.Rows
	if len(keywords) != 0 {
		query := `SELECT (id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes, rating, views) 
			FROM
				((SELECT id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes, rating, views
				FROM forum.messages
				WHERE message_ts @@ to_tsquery_multilang($1))
			UNION
				(SELECT id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes, rating, views
				FROM forum.messages
				WHERE tags && ($2) COLLATE case_insensitive))`

		switch order {
		case "newest-updated":
			query = query + ` ORDER BY updated_at DESC`
		case "highest-rated":
			query = query + ` ORDER BY rating DESC`
		case "most-views":
			query = query + ` ORDER BY views DESC`
		default:
			query = query + ` ORDER BY updated_at DESC`
		}
		if limit != 0 {
			query = query + fmt.Sprintf(` LIMIT %v`, limit)
		}
		rows, err = conn.Query(context.Background(),
			query, strings.Join(keywords, " | "), keywords,
		)
		if err != nil {
			return nil, err
		}
	} else {
		query := `SELECT (id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes, rating, views)
			FROM forum.messages`

		switch order {
		case "newest-updated":
			query = query + ` ORDER BY updated_at DESC`
		case "highest-rated":
			query = query + ` ORDER BY rating DESC`
		case "most-views":
			query = query + ` ORDER BY views DESC`
		default:
			query = query + ` ORDER BY updated_at DESC`
		}
		if limit != 0 {
			query = query + fmt.Sprintf(` LIMIT %v`, limit)
		}
		rows, err = conn.Query(context.Background(), query)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, r.NotFoundErr()
	}
	return &messages, nil
}

func (r *PsqlForumRepository) FindMessagesByThreadID(thread_id int) (*[]models.Message, error) {
	var messages []models.Message

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`UPDATE forum.threads SET 
		views=views+1
		WHERE id = $1`, thread_id,
	)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(context.Background(),
		`SELECT (id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes, rating, views)
		FROM forum.messages WHERE thread_id=$1`,
		thread_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, r.NotFoundErr()
	}
	return &messages, nil
}

func (r *PsqlForumRepository) UpdateMessage(id int, message models.Message) (*models.Message, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE forum.messages SET 
		thread_id=COALESCE($1, thread_id), user_id=COALESCE($2, user_id), title=COALESCE($3, title), 
		body=COALESCE($4, body), tags=COALESCE($5, tags), updated_at=NOW(),
		upvotes=COALESCE($6, upvotes), downvotes=COALESCE($7, downvotes)
		WHERE id = $8
		RETURNING
		(id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes, rating, views)`,
		message.ThreadID, message.UserID, message.Title, message.Body,
		message.Tags, message.Upvotes, message.Downvotes, id,
	).Scan(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *PsqlForumRepository) DeleteMessage(id int) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.messages WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	_, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("messages/%v", id), "PATCH")
	if err != nil {
		return err
	}
	_, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("messages/%v", id), "DELETE")
	if err != nil {
		return err
	}
	return nil
}

func (r *PsqlForumRepository) DeleteMessagesOfThread(threadID int) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT id FROM forum.messages WHERE thread_id = $1`, threadID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var message int
		err = rows.Scan(&message)
		if err != nil {
			return err
		}
		_, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("messages/%v", message), "PATCH")
		if err != nil {
			return err
		}
		_, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("messages/%v", message), "DELETE")
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
