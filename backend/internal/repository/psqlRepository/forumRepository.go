package psqlRepository

import (
	"context"
	"gamehangar/internal/domain/models"
)

type PsqlForumRepository struct {
	databaseClient psqlDatabaseClient
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlForumRepository(dbClient psqlDatabaseClient) *PsqlForumRepository {
	return &PsqlForumRepository{
		databaseClient: dbClient,
	}
}

func (r *PsqlForumRepository) NotFoundErr() error { return r.databaseClient.ErrNoRows() }

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
		(id, name)`,
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
		`SELECT (id, name) FROM forum.topics WHERE id = $1 LIMIT 1`,
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

	rows, err := conn.Query(context.Background(), `SELECT (id, name) FROM forum.topics`)
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
	return &topics, nil
}

func (r *PsqlForumRepository) UpdateTopic(id int, topic models.Topic) (*models.Topic, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE forum.topics SET 
		name=COALESCE($1, name)
		WHERE id = $2
		RETURNING
		(id, name)`,
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
			(title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			(id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes)`,
		thread.Title, thread.UserID, thread.TopicID, thread.Tags,
		thread.CreatedAt, thread.LastUpdate, thread.TotalUpvotes, thread.TotalDownvotes,
	).Scan(&thread)
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
		`SELECT (id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes)
		FROM forum.threads WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&thread)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r *PsqlForumRepository) FindThreads() (*[]models.Thread, error) {
	var threads []models.Thread

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT (id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes)
		FROM forum.threads`,
	)
	if err != nil {
		return nil, err
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
			title=COALESCE($1, title), user_id=COALESCE($2, user_id), 
		topic_id=COALESCE($3, topic_id), tags=COALESCE($4, tags), 
			created_at=COALESCE($5, created_at), updated_at=COALESCE($6, updated_at), 
		upvotes=COALESCE($7, upvotes), downvotes=COALESCE($8, downvotes)
			WHERE id = $9
		RETURNING
			(id, title, user_id, topic_id, tags, created_at, updated_at, upvotes, downvotes)`,
		thread.Title, thread.UserID, thread.TopicID, thread.Tags, thread.CreatedAt,
		thread.LastUpdate, thread.TotalUpvotes, thread.TotalDownvotes, id,
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
		(thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING
		(id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes)`,
		message.ThreadID, message.UserID, message.Title, message.Body,
		message.Tags, message.CreatedAt, message.UpdatedAt, message.Upvotes,
		message.Downvotes,
	).Scan(&message)
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
		`SELECT (id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes)
		FROM forum.messages WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *PsqlForumRepository) FindMessages() (*[]models.Message, error) {
	var messages []models.Message

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT (id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes)
		FROM forum.messages`,
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
	return &messages, nil
}

func (r *PsqlForumRepository) FindMessagesByThreadID(thread_id int) (*[]models.Message, error) {
	var messages []models.Message

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT (id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes)
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
		body=COALESCE($4, body), tags=COALESCE($5, tags), created_at=COALESCE($6, created_at),
		updated_at=COALESCE($7, updated_at), upvotes=COALESCE($8, upvotes), downvotes=COALESCE($9, downvotes)
		WHERE id = $10
		RETURNING
		(id, thread_id, user_id, title, body, tags, created_at, updated_at, upvotes, downvotes)`,
		message.ThreadID, message.UserID, message.Title, message.Body,
		message.Tags, message.CreatedAt, message.UpdatedAt, message.Upvotes,
		message.Downvotes, id,
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
	return nil
}
