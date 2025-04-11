package psqlRepository

import (
	"context"
	"errors"

	"gamehangar/internal/domain/models"
)

type PsqlDemoRepository struct {
	databaseClient psqlDatabaseClient
	notFoundErr    error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlDemoRepository(dbClient psqlDatabaseClient) (*PsqlDemoRepository, error) {
	return &PsqlDemoRepository{
		databaseClient: dbClient,
		notFoundErr:    errors.New("Not Found"),
	}, nil
}

func (r *PsqlDemoRepository) CreateDemo(demo models.Demo) (*models.Demo, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO demo.demos
		(id, name, description, link, userID, createdAt, updatedAt, upvotes, downvotes, threadID) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
		(id, name, description, link, userID, createdAt, updatedAt, upvotes, downvotes, threadID)`,
		demo.ID, demo.Name, demo.Description, demo.Link, demo.UserID,
		demo.CreatedAt, demo.UpdatedAt, demo.Upvotes, demo.Downvotes, demo.ThreadID,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

func (r *PsqlDemoRepository) FindDemoByID(id string) (*models.Demo, error) {
	var demo models.Demo
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM demo.demos WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&demo.ID, &demo.Name, &demo.Description, &demo.Link, &demo.UserID,
		&demo.CreatedAt, &demo.UpdatedAt, &demo.Upvotes, &demo.Downvotes, &demo.ThreadID)
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

func (r *PsqlDemoRepository) FindDemos() (*[]models.Demo, error) {
	var demos []models.Demo

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `SELECT * FROM demo.demos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var demo models.Demo
		err = rows.Scan(&demo.ID, &demo.Name, &demo.Description, &demo.Link, &demo.UserID,
			&demo.CreatedAt, &demo.UpdatedAt, &demo.Upvotes, &demo.Downvotes, &demo.ThreadID)
		if err != nil {
			return nil, err
		}
		demos = append(demos, demo)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &demos, nil
}

func (r *PsqlDemoRepository) UpdateDemo(id string, demo models.Demo) (*models.Demo, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE demo.demos SET 
		name=$1, description=$2, link=$3, userID=$4, createdAt=$5, updatedAt=$6, upvotes=$7, downvotes=$8, threadID=$9 
		WHERE id = $10
		RETURNING
		(id, name, description, link, userID, createdAt, updatedAt, upvotes, downvotes, threadID)`,
		demo.Name, demo.Description, demo.Link, demo.UserID, demo.CreatedAt,
		demo.UpdatedAt, demo.Upvotes, demo.Downvotes, demo.ThreadID, id,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}
	return &demo, err
}

func (r *PsqlDemoRepository) DeleteDemo(id string) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM demo.demos WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return r.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}
