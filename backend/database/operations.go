package database

import (
	"context"
)

func CreateDemo(demo Demo) (string, error) {
	conn, err := dbpool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		return "", err
	}

	row := conn.QueryRow(context.Background(),
		`INSERT INTO demos
			(id, name, description, link, created_at, updated_at, upvotes, downvotes) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`,
		demo.ID, demo.Name, demo.Description, demo.Link,
		demo.Created_at, demo.Updated_at, demo.Upvotes, demo.Downvotes,
	)

	var id string
	err = row.Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
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
	).Scan(&demo.ID, &demo.Name, &demo.Description, &demo.Link,
		&demo.Created_at, &demo.Updated_at, &demo.Upvotes, &demo.Downvotes)
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

// func FindDemos(w http.ResponseWriter, r *http.Request) (*[]Demo, bool) {
// 	return &[]demos, true
// }
