package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// TODO: MapNameToUser + assert one unique user per username!
// TODO: A new demo creates a new topic on forums — a different service!

func CreateDemo(demo Demo) (*Demo, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO demos
			(id, name, description, link, user_id, created_at, updated_at, upvotes, downvotes, topic_id) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
			(id, name, description, link, user_id, created_at, updated_at, upvotes, downvotes, topic_id)`,
		demo.ID, demo.Name, demo.Description, demo.Link, demo.User_id,
		demo.Created_at, demo.Updated_at, demo.Upvotes, demo.Downvotes, demo.Topic_id,
	)

	err = row.Scan(&demo)
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

func FindFirstDemo(id string) (*Demo, error) {
	var demo Demo
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM demos WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&demo.ID, &demo.Name, &demo.Description, &demo.Link, &demo.User_id,
		&demo.Created_at, &demo.Updated_at, &demo.Upvotes, &demo.Downvotes, &demo.Topic_id)
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

func FindDemos() (*[]Demo, error) {
	var demos []Demo

	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT * FROM demos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var demo Demo
		demos = append(demos, demo)
		err = rows.Scan(&demo.ID, &demo.Name, &demo.Description, &demo.Link, &demo.User_id,
			&demo.Created_at, &demo.Updated_at, &demo.Upvotes, &demo.Downvotes, &demo.Topic_id)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &demos, nil
}

func UpdateDemo(demo Demo) (*Demo, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE demos SET 
		name=$1, description=$2, link=$3, user_id=$4, created_at=$5, updated_at=$6, upvotes=$7, downvotes=$8, topic_id=$9 
		WHERE id=$10`,
		demo.Name, demo.Description, demo.Link, demo.User_id, demo.Created_at,
		demo.Updated_at, demo.Upvotes, demo.Downvotes, demo.Topic_id, demo.ID,
	)
	if ct.RowsAffected() == 0 {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

func DeleteDemo(id string) error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM demos WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	return nil
}
