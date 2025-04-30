package psqlRepository

import (
	"context"
	"strings"

	"gamehangar/internal/domain/models"
)

type PsqlDemoRepository struct {
	databaseClient psqlDatabaseClient
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlDemoRepository(dbClient psqlDatabaseClient) *PsqlDemoRepository {
	return &PsqlDemoRepository{
		databaseClient: dbClient,
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
		(title, description, link, tags, user_id, thread_id) 
		VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING
		(id, title, description, link, tags, user_id, thread_id)`,
		demo.Title, demo.Description, demo.Link, demo.Tags, demo.UserID, demo.ThreadID,
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
		`SELECT (id, title, description, link, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes)
		FROM demo.demos WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

func (r *PsqlDemoRepository) FindDemosByQuery(query *[]string) (*[]models.Demo, error) {
	var demos []models.Demo

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`(
		SELECT (id, title, description, link, tags, user_id,
		thread_id, created_at, updated_at, upvotes, downvotes)
		FROM demo.demos
		WHERE demo_ts @@ to_tsquery_multilang($1)
		)
		UNION
		(
		SELECT (id, title, description, link, tags, user_id,
		thread_id, created_at, updated_at, upvotes, downvotes)
		FROM demo.demos
		WHERE tags && ($2) COLLATE case_insensitive
		ORDER BY updated_at DESC
		)`, strings.Join(*query, " | "), *query,
	)
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
	if len(demos) == 0 {
		return nil, r.NotFoundErr()
	}
	return &demos, nil
}

func (r *PsqlDemoRepository) FindDemos() (*[]models.Demo, error) {
	var demos []models.Demo

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT (id, title, description, link, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes) 
		FROM demo.demos`,
	)
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
		link=COALESCE($3, link), tags=COALESCE($4, tags), user_id=COALESCE($5, user_id),
			thread_id=COALESCE($6, thread_id), updated_at=NOW(),
		upvotes=COALESCE($7, upvotes), downvotes=COALESCE($8, downvotes)
			WHERE id = $9
		RETURNING
			(id, title, description, link, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes)`,
		demo.Title, demo.Description, demo.Link, demo.Tags, demo.UserID, demo.ThreadID,
		demo.Upvotes, demo.Downvotes, id,
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
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	return nil
}
