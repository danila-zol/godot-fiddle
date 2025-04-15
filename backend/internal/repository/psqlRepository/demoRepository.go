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
func NewPsqlDemoRepository(dbClient psqlDatabaseClient) *PsqlDemoRepository {
	return &PsqlDemoRepository{
		databaseClient: dbClient,
		notFoundErr:    errors.New("Not Found"),
	}
}

func (r *PsqlDemoRepository) NotFoundErr() error { return r.databaseClient.ErrNoRows() }

func (r *PsqlDemoRepository) CreateDemo(demo models.Demo) (*models.Demo, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO demo.demos
		(title, description, link, tags, "userID", "threadID", "createdAt", "updatedAt", upvotes, downvotes) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
		(id, title, description, link, tags, "userID", "threadID", "createdAt", "updatedAt", upvotes, downvotes)`,
		demo.Title, demo.Description, demo.Link, demo.Tags, demo.UserID,
		demo.ThreadID, demo.CreatedAt, demo.UpdatedAt, demo.Upvotes, demo.Downvotes,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

func (r *PsqlDemoRepository) FindDemoByID(id int) (*models.Demo, error) {
	var demo models.Demo
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM demo.demos WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&demo)
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
		err = rows.Scan(&demo)
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

func (r *PsqlDemoRepository) UpdateDemo(id int, demo models.Demo) (*models.Demo, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE demo.demos SET 
			title=COALESCE($1, title), description=COALESCE($2, description),
		link=COALESCE($3, link), tags=COALESCE($4, tags), "userID"=COALESCE($5, "userID"),
			"threadID"=COALESCE($6, "threadID"), "createdAt"=COALESCE($7, "createdAt"),
		"updatedAt"=COALESCE($8, "updatedAt"), upvotes=COALESCE($9, upvotes), downvotes=COALESCE($10, downvotes)
			WHERE id = $11
		RETURNING
			(id, title, description, link, tags "userID", "threadID", "createdAt", "updatedAt", upvotes, downvotes)`,
		demo.Title, demo.Description, demo.Link, demo.Tags, demo.UserID, demo.ThreadID,
		demo.CreatedAt, demo.UpdatedAt, demo.Upvotes, demo.Downvotes, id,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}
	return &demo, err
}

func (r *PsqlDemoRepository) DeleteDemo(id int) error {
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
