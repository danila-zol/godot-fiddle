package psqlRepository

import (
	"context"
	"errors"
	"gamehangar/internal/database"
	"gamehangar/internal/domain/models"
)

type PsqlForumRepository struct {
	databaseClient *database.PsqlDatabaseClient
	notFoundErr    error
}

func NewPsqlForumRepository(dbClient *database.PsqlDatabaseClient) (*PsqlForumRepository, error) {
	return &PsqlForumRepository{
		databaseClient: dbClient,
		notFoundErr:    errors.New("Not Found"),
	}, nil
}

func (pfr *PsqlForumRepository) CreateTopic(topic models.Topic) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`INSERT INTO forum.topics
		(id, name) 
		VALUES
		($1, $2)`,
		topic.ID, topic.Name)

	if err != nil {
		return err
	}
	return nil
}

func (pfr *PsqlForumRepository) FindByTopicID(id string) (*models.Topic, error) {
	var topic models.Topic
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
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

func (pfr *PsqlForumRepository) FindTopics() (*[]models.Topic, error) {
	var topics []models.Topic

	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
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

func (pfr *PsqlForumRepository) UpdateTopic(topic models.Topic) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE forum.topics SET 
		name=$1
		WHERE id = $2`,
		topic.Name, topic.ID,
	)
	if ct.RowsAffected() == 0 {
		return pfr.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}

func (pfr *PsqlForumRepository) DeleteTopic(id string) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.topics WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pfr.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}

func (pfr *PsqlForumRepository) CreateThread(thread models.Thread) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`INSERT INTO forum.threads
		(id, title, userID, topicID, tags, createdAt, lastUpdate, totalUpvotes, totalDownvotes) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		thread.ID, thread.Title, thread.UserID, thread.TopicID, thread.Tags,
		thread.CreatedAt, thread.LastUpdate, thread.TotalUpvotes, thread.TotalDownvotes,
	)
	if err != nil {
		return err
	}

	return nil
}

func (pfr *PsqlForumRepository) FindByThreadID(id string) (*models.Thread, error) {
	var thread models.Thread
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
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

func (pfr *PsqlForumRepository) FindThreads() (*[]models.Thread, error) {
	var threads []models.Thread

	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
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

func (pfr *PsqlForumRepository) UpdateThread(id string, thread models.Thread) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE forum.threads SET 
		name=$1, userID=$2, topicID=$3, tag=$4, createdAt=$5, lastUpdate=$6, totalUpvotes=$7, totalDownvotes=$8
		WHERE id = $9`,
		thread.Title, thread.UserID, thread.TopicID, thread.Tags, thread.CreatedAt,
		thread.LastUpdate, thread.TotalUpvotes, thread.TotalDownvotes, thread.ID,
	)
	if ct.RowsAffected() == 0 {
		return pfr.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}

func (pfr *PsqlForumRepository) DeleteThread(id string) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
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

func (pfr *PsqlForumRepository) CreateMessage(message models.Message) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`INSERT INTO forum.messages
		(id, threadID, userID, title, body, tags, createdAt, updatedAt, upvotes, downvotes) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		message.ID, message.ThreadID, message.UserID, message.Title, message.Body,
		message.Tags, message.CreatedAt, message.UpdatedAt, message.Upvotes,
		message.Downvotes,
	)

	if err != nil {
		return err
	}
	return nil
}

func (pfr *PsqlForumRepository) FindByMessageID(id string) (*models.Message, error) {
	var message models.Message
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
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

func (pfr *PsqlForumRepository) FindMessages() (*[]models.Message, error) {
	var messages []models.Message

	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
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

func (pfr *PsqlForumRepository) FindMessagesByThreadID(threadID string) (*[]models.Message, error) {
	var messages []models.Message

	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
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

func (pfr *PsqlForumRepository) UpdateMessage(id string, message models.Message) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE forum.messages SET 
		threadID=$1, userID=$2, title=$3, body=$4, tags=$5, createdAt=$6, updatedAt=$7, upvotes=$8, downvotes=$9
		WHERE id = $10`,
		message.ThreadID, message.UserID, message.Title, message.Body,
		message.Tags, message.CreatedAt, message.UpdatedAt, message.Upvotes,
		message.Downvotes, message.ID,
	)
	if ct.RowsAffected() == 0 {
		return pfr.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}

func (pfr *PsqlForumRepository) DeleteMessage(id string) error {
	conn, err := pfr.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.messages WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pfr.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}
