package psqlRepository

import (
	"context"
	"fmt"
	"strings"

	"gamehangar/internal/domain/models"

	"github.com/jackc/pgx/v5"
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
		(id, title, description, link, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes, views)`,
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
		`UPDATE demo.demos SET views=views+1
		WHERE id = $1
		RETURNING
		(id, title, description, link, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes, views)`,
		id,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}
	return &demo, nil
}

func (r *PsqlDemoRepository) FindDemos(keywords []string, limit uint64) (*[]models.Demo, error) {
	var demos []models.Demo

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var rows pgx.Rows
	if len(keywords) != 0 {
		query := `SELECT (id, title, description, link, tags, user_id,
			thread_id, created_at, updated_at, upvotes, downvotes, views)
		FROM 
			((SELECT id, title, description, link, tags, user_id,
				thread_id, created_at, updated_at, upvotes, downvotes, views
			FROM demo.demos
			WHERE demo_ts @@ to_tsquery_multilang($1))
			UNION
			(SELECT id, title, description, link, tags, user_id,
				thread_id, created_at, updated_at, upvotes, downvotes, views
			FROM demo.demos
			WHERE tags && ($2) COLLATE case_insensitive))
		ORDER BY updated_at DESC`
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
		query := `SELECT (id, title, description, link, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes, views) 
		FROM demo.demos`
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
			(id, title, description, link, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes, views)`,
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
