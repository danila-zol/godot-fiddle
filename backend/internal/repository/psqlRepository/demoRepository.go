package psqlRepository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"gamehangar/internal/domain/models"

	"github.com/jackc/pgx/v5"
)

type PsqlDemoRepository struct {
	databaseClient       psqlDatabaseClient
	objectUploader       ObjectUploader
	enforcer             Enforcer
	attachmentMissingErr error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlDemoRepository(dbClient psqlDatabaseClient, o ObjectUploader, e Enforcer) *PsqlDemoRepository {
	return &PsqlDemoRepository{
		databaseClient:       dbClient,
		attachmentMissingErr: errors.New("Missing attachment!"),
		objectUploader:       o,
		enforcer:             e,
	}
}

func (r *PsqlDemoRepository) NotFoundErr() error { return r.databaseClient.ErrNoRows() }

func (r *PsqlDemoRepository) CreateDemo(demo models.Demo, demoFile, demoThumbnail io.Reader) (*models.Demo, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO demo.demos
		(title, description, tags, user_id, thread_id) 
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING
		(id, title, description, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes, rating, views, object_key, thumbnail_key)`,
		demo.Title, demo.Description, demo.Tags, demo.UserID, demo.ThreadID,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}
	_, err = r.enforcer.AddPermissions(demo.UserID.String(), fmt.Sprintf("demos/%v", *demo.ID), "PATCH")
	if err != nil {
		return nil, err
	}
	_, err = r.enforcer.AddPermissions(demo.UserID.String(), fmt.Sprintf("demos/%v", *demo.ID), "DELETE")
	if err != nil {
		return nil, err
	}

	if demoFile == nil {
		return nil, r.attachmentMissingErr
	}
	err = r.objectUploader.PutObject(*demo.Key, demoFile)
	if err != nil {
		return nil, err
	}
	demo.Key, err = r.objectUploader.GetObjectLink(*demo.Key)
	if err != nil {
		return nil, err
	}

	if demoThumbnail == nil {
		return nil, r.attachmentMissingErr
	}
	err = r.objectUploader.PutObject(*demo.ThumbnailKey, demoThumbnail)
	if err != nil {
		return nil, err
	}
	demo.ThumbnailKey, err = r.objectUploader.GetObjectLink(*demo.ThumbnailKey)
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
		(id, title, description, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes, rating, views, object_key, thumbnail_key)`,
		id,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}

	demo.Key, err = r.objectUploader.GetObjectLink(*demo.Key)
	if err != nil {
		return nil, err
	}
	demo.ThumbnailKey, err = r.objectUploader.GetObjectLink(*demo.ThumbnailKey)
	if err != nil {
		return nil, err
	}

	return &demo, nil
}

func (r *PsqlDemoRepository) FindDemos(keywords []string, limit uint64, order string) (*[]models.Demo, error) {
	var demos []models.Demo

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var rows pgx.Rows
	if len(keywords) != 0 {
		query := `SELECT (id, title, description, tags, user_id,
			thread_id, created_at, updated_at, upvotes, downvotes, rating, views, object_key, thumbnail_key)
			FROM 
			((SELECT id, title, description, tags, user_id,
				thread_id, created_at, updated_at, upvotes, downvotes, rating, views, object_key, thumbnail_key
			FROM demo.demos
			WHERE demo_ts @@ to_tsquery_multilang($1))
			UNION
			(SELECT id, title, description, tags, user_id,
				thread_id, created_at, updated_at, upvotes, downvotes, rating, views, object_key, thumbnail_key
			FROM demo.demos
			WHERE tags && ($2) COLLATE case_insensitive))`

		switch order {
		case "newest-updated":
			query = query + ` ORDER BY updated_at DESC`
		case "highest-rated":
			query = query + ` ORDER BY rating DESC`
		case "most-views":
			query = query + ` ORDER BY views DESC`
		default:
			query = query + ` ORDER BY updated_at DESC`
		}
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
		query := `SELECT (id, title, description, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes, rating, views, object_key, thumbnail_key) 
		FROM demo.demos`

		switch order {
		case "newest-updated":
			query = query + ` ORDER BY updated_at DESC`
		case "highest-rated":
			query = query + ` ORDER BY upvotes DESC`
		case "most-views":
			query = query + ` ORDER BY views DESC`
		default:
			query = query + ` ORDER BY updated_at DESC`
		}
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
		demo.Key, _ = r.objectUploader.GetObjectLink(*demo.Key)
		demo.ThumbnailKey, _ = r.objectUploader.GetObjectLink(*demo.ThumbnailKey)
		demos = append(demos, demo)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &demos, nil
}

func (r *PsqlDemoRepository) UpdateDemo(id int, demo models.Demo, demoFile, demoThumbnail io.Reader) (*models.Demo, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE demo.demos SET 
			title=COALESCE($1, title), description=COALESCE($2, description),
		tags=COALESCE($3, tags), user_id=COALESCE($4, user_id),
			thread_id=COALESCE($5, thread_id), updated_at=NOW(),
		upvotes=COALESCE($6, upvotes), downvotes=COALESCE($7, downvotes)
			WHERE id = $8
		RETURNING
			(id, title, description, tags, user_id, thread_id, created_at, updated_at, upvotes, downvotes, rating, views, object_key, thumbnail_key)`,
		demo.Title, demo.Description, demo.Tags, demo.UserID, demo.ThreadID,
		demo.Upvotes, demo.Downvotes, id,
	).Scan(&demo)
	if err != nil {
		return nil, err
	}

	if demoFile != nil {
		err = r.objectUploader.PutObject(*demo.Key, demoFile)
		if err != nil {
			return nil, err
		}
	}

	if demoThumbnail != nil {
		err = r.objectUploader.PutObject(*demo.ThumbnailKey, demoThumbnail)
		if err != nil {
			return nil, err
		}
	}
	demo.Key, _ = r.objectUploader.GetObjectLink(*demo.Key)
	demo.ThumbnailKey, _ = r.objectUploader.GetObjectLink(*demo.ThumbnailKey)

	return &demo, err
}

func (r *PsqlDemoRepository) DeleteDemo(id int) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	err = r.objectUploader.DeleteObject(fmt.Sprintf("demo-%v", id))
	if err != nil {
		return err
	}
	err = r.objectUploader.DeleteObject(fmt.Sprintf("demo-thumbnail-%v", id))
	if err != nil {
		return err
	}

	ct, err := conn.Exec(context.Background(), `DELETE FROM demo.demos WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	changed, err := r.enforcer.RemovePermissionsForObject(fmt.Sprintf("demos/%v", id), "PATCH")
	if err != nil || !changed {
		return errors.New("Somehow did not changed anything WTF?")
	}
	changed, err = r.enforcer.RemovePermissionsForObject(fmt.Sprintf("demos/%v", id), "DELETE")
	if err != nil || !changed {
		return err
	}
	return nil
}
