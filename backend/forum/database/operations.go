package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTopic(topic Topic) (*Topic, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO forum.topics
			(id, name) 
		VALUES
			($1, $2)
		RETURNING
			(id, name)`,
		topic.ID, topic.Name)

	err = row.Scan(&topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func FindFirstTopic(id string) (*Topic, error) {
	var topic Topic
	conn, err := dbpool.Acquire(context.Background())
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

func FindTopics() (*[]Topic, error) {
	var topics []Topic

	conn, err := dbpool.Acquire(context.Background())
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
		var topic Topic
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

func UpdateTopic(topic Topic) (*Topic, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE forum.topics SET 
		name=$1
		WHERE id = $2`,
		topic.Name, topic.ID,
	)
	if ct.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func DeleteTopic(id string) error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.topics WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	return nil
}

func CreateThread(thread Thread) (*Thread, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO forum.threads
		(id, title, user_id, topic_id, tag, created_at, last_update, total_upvotes, total_downvotes) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING
		(id, title, user_id, topic_id, tag, created_at, last_update, total_upvotes, total_downvotes)`,
		thread.ID, thread.Title, thread.User_id, thread.Topic_id, thread.Tag,
		thread.Created_at, thread.Last_update, thread.Total_upvotes, thread.Total_downvotes,
	)

	err = row.Scan(&thread)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func FindFirstThread(id string) (*Thread, error) {
	var thread Thread
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM forum.threads WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&thread.ID, &thread.Title, &thread.User_id, &thread.Topic_id, &thread.Tag,
		&thread.Created_at, &thread.Last_update, &thread.Total_upvotes, &thread.Total_downvotes)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func FindThreads() (*[]Thread, error) {
	var threads []Thread

	conn, err := dbpool.Acquire(context.Background())
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
		var thread Thread
		err = rows.Scan(&thread.ID, &thread.Title, &thread.User_id, &thread.Topic_id,
			&thread.Tag, &thread.Created_at, &thread.Last_update,
			&thread.Total_upvotes, &thread.Total_downvotes)
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

func UpdateThread(thread Thread) (*Thread, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE forum.threads SET 
		name=$1, user_id=$2, topic_id=$3, tag=$4, created_at=$5, last_update=$6, total_upvotes=$7, total_downvotes=$8
		WHERE id = $9`,
		thread.Title, thread.User_id, thread.Topic_id, thread.Tag, thread.Created_at,
		thread.Last_update, thread.Total_upvotes, thread.Total_downvotes, thread.ID,
	)
	if ct.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func DeleteThread(id string) error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.threads WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	return nil
}

func CreateMessage(message Message) (*Message, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO forum.messages
		(id, thread_id, user_id, title, body, tag, created_at, updated_at, upvotes, downvotes) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
		(id, thread_id, user_id, title, body, tag, created_at, updated_at, upvotes, downvotes)`,
		message.ID, message.Thread_id, message.User_id, message.Title, message.Body,
		message.Tag, message.Created_at, message.Updated_at, message.Upvotes,
		message.Downvotes,
	)

	err = row.Scan(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func FindFirstMessage(id string) (*Message, error) {
	var message Message
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM forum.messages WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&message.ID, &message.Thread_id, &message.User_id, &message.Title, &message.Body,
		&message.Tag, &message.Created_at, &message.Updated_at, &message.Upvotes,
		&message.Downvotes)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func FindMessages() (*[]Message, error) {
	var messages []Message

	conn, err := dbpool.Acquire(context.Background())
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
		var message Message
		err = rows.Scan(&message.ID, &message.Thread_id, &message.User_id, &message.Title,
			&message.Body, &message.Tag, &message.Created_at, &message.Updated_at,
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

func UpdateMessage(message Message) (*Message, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE forum.messages SET 
		thread_id=$1, user_id=$2, title=$3, body=$4, tag=$5, created_at=$6, updated_at=$7, upvotes=$8, downvotes=$9
		WHERE id = $10`,
		message.Thread_id, message.User_id, message.Title, message.Body,
		message.Tag, message.Created_at, message.Updated_at, message.Upvotes,
		message.Downvotes, message.ID,
	)
	if ct.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func DeleteMessage(id string) error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM forum.messages WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	return nil
}
