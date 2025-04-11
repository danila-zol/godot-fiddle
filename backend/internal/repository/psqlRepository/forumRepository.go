package psqlRepository

import (
	"context"
	"errors"
	"gamehangar/internal/domain/models"
)

type PsqlForumRepository struct {
	databaseClient psqlDatabaseClient
	notFoundErr    error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlForumRepository(dbClient psqlDatabaseClient) *PsqlForumRepository {
	return &PsqlForumRepository{
		databaseClient: dbClient,
		notFoundErr:    errors.New("Not Found"),
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
		(id, name) 
		VALUES
		($1, $2)
		RETURNING
		(id, name)`,
		topic.ID, topic.Name,
	).Scan(&topic)

	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *PsqlForumRepository) FindTopicByID(id string) (*models.Topic, error) {
	var topic models.Topic
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM forum.topics WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&topic.ID, &topic.Name)
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

	rows, err := conn.Query(context.Background(), `SELECT * FROM forum.topics`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var topic models.Topic
		err = rows.Scan(&topic.ID, &topic.Name)
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

func (r *PsqlForumRepository) UpdateTopic(id string, topic models.Topic) (*models.Topic, error) {
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

func (r *PsqlForumRepository) DeleteTopic(id string) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.topics WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return r.notFoundErr
	}
	if err != nil {
		return err
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
		(id, title, "userID", "topicID", tags, "createdAt", "lastUpdate", "totalUpvotes", "totalDownvotes") 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING
		(id, title, "userID", "topicID", tags, "createdAt", "lastUpdate", "totalUpvotes", "totalDownvotes")`,
		thread.ID, thread.Title, thread.UserID, thread.TopicID, thread.Tags,
		thread.CreatedAt, thread.LastUpdate, thread.TotalUpvotes, thread.TotalDownvotes,
	).Scan(&thread)
	if err != nil {
		return nil, err
	}

	return &thread, err
}

func (r *PsqlForumRepository) FindThreadByID(id string) (*models.Thread, error) {
	var thread models.Thread
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM forum.threads WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&thread.ID, &thread.Title, &thread.UserID, &thread.TopicID, &thread.Tags,
		&thread.CreatedAt, &thread.LastUpdate, &thread.TotalUpvotes, &thread.TotalDownvotes)
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

	rows, err := conn.Query(context.Background(), `SELECT * FROM forum.threads`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var thread models.Thread
		err = rows.Scan(&thread.ID, &thread.Title, &thread.UserID, &thread.TopicID,
			&thread.Tags, &thread.CreatedAt, &thread.LastUpdate,
			&thread.TotalUpvotes, &thread.TotalDownvotes)
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

func (r *PsqlForumRepository) UpdateThread(id string, thread models.Thread) (*models.Thread, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE forum.threads SET 
		name=COALESCE($1, name), "userID"=COALESCE($2, "userID"), "topicID"=COALESCE($3, "topicID"), tags=COALESCE($4, tags), "createdAt"=COALESCE($5, "createdAt"), 
		"lastUpdate"=COALESCE($6, "lastUpdate"), "totalUpvotes"=COALESCE($7, "totalUpvotes"), "totalDownvotes"=COALESCE($8, "totalDownvotes")
		WHERE id = $9
		RETURNING
		(id, title, "userID", "topicID", tags, "createdAt", "lastUpdate", "totalUpvotes", "totalDownvotes")`,
		thread.Title, thread.UserID, thread.TopicID, thread.Tags, thread.CreatedAt,
		thread.LastUpdate, thread.TotalUpvotes, thread.TotalDownvotes, id,
	).Scan(&thread)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r *PsqlForumRepository) DeleteThread(id string) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.threads WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return errors.New("Not Found")
	}
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
		(id, "threadID", "userID", title, body, tags, "createdAt", "updatedAt", upvotes, downvotes) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
		(id, "threadID", "userID", title, body, tags, "createdAt", "updatedAt", upvotes, downvotes)`,
		message.ID, message.ThreadID, message.UserID, message.Title, message.Body,
		message.Tags, message.CreatedAt, message.UpdatedAt, message.Upvotes,
		message.Downvotes,
	).Scan(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *PsqlForumRepository) FindMessageByID(id string) (*models.Message, error) {
	var message models.Message
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM forum.messages WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&message.ID, &message.ThreadID, &message.UserID, &message.Title, &message.Body,
		&message.Tags, &message.CreatedAt, &message.UpdatedAt, &message.Upvotes,
		&message.Downvotes)
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

	rows, err := conn.Query(context.Background(), `SELECT * FROM forum.messages`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message.ID, &message.ThreadID, &message.UserID, &message.Title,
			&message.Body, &message.Tags, &message.CreatedAt, &message.UpdatedAt,
			&message.Upvotes, &message.Downvotes)
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

func (r *PsqlForumRepository) FindMessagesByThreadID(threadID string) (*[]models.Message, error) {
	var messages []models.Message

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT * FROM forum.messages WHERE threadID=$1`,
		threadID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message.ID, &message.ThreadID, &message.UserID, &message.Title,
			&message.Body, &message.Tags, &message.CreatedAt, &message.UpdatedAt,
			&message.Upvotes, &message.Downvotes)
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

func (r *PsqlForumRepository) UpdateMessage(id string, message models.Message) (*models.Message, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE forum.messages SET 
		"threadID"=COALESCE($1, "threadID"), "userID"=COALESCE($2, "userID"), title=COALESCE($3, title), 
		body=COALESCE($4, body), tags=COALESCE($5, tags), "createdAt"=COALESCE($6, "createdAt"),
		"updatedAt"=COALESCE($7, "updatedAt"), upvotes=COALESCE($8, upvotes), downvotes=COALESCE($9, downvotes)
		WHERE id = $10
		RETURNING
		(id, "threadID", "userID", title, body, tags, "createdAt", "updatedAt", upvotes, downvotes)`,
		message.ThreadID, message.UserID, message.Title, message.Body,
		message.Tags, message.CreatedAt, message.UpdatedAt, message.Upvotes,
		message.Downvotes, id,
	).Scan(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *PsqlForumRepository) DeleteMessage(id string) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.messages WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return r.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}
