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

type PsqlAssetRepository struct {
	databaseClient       psqlDatabaseClient
	objectUploader       ObjectUploader
	attachmentMissingErr error
	conflictErr          error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlAssetRepository(dbClient psqlDatabaseClient, o ObjectUploader) *PsqlAssetRepository {
	return &PsqlAssetRepository{
		databaseClient:       dbClient,
		attachmentMissingErr: errors.New("Missing attachment!"),
		conflictErr:          errors.New("Record conflict!"),
		objectUploader:       o,
	}
}

func (r *PsqlAssetRepository) NotFoundErr() error { return r.databaseClient.ErrNoRows() }

// Returns "Record conflict!" to specify conflicting record versions on update
func (r *PsqlAssetRepository) ConflictErr() error { return r.conflictErr }

func (r *PsqlAssetRepository) CreateAsset(asset models.Asset, assetFile, assetThumbnail io.Reader) (*models.Asset, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO asset.assets
		(name, description, tags) 
		VALUES
		($1, $2, $3)
		RETURNING
		(id, name, description, tags, created_at, updated_at, version, upvotes, downvotes, rating, views, object_key, thumbnail_key)`,
		asset.Name, asset.Description, asset.Tags,
	).Scan(&asset)
	if err != nil {
		return nil, err
	}

	if assetFile == nil {
		return nil, r.attachmentMissingErr
	}
	err = r.objectUploader.PutObject(*asset.Key, assetFile)
	if err != nil {
		return nil, err
	}
	asset.Key, err = r.objectUploader.GetObjectLink(*asset.Key)
	if err != nil {
		return nil, err
	}

	if assetThumbnail == nil {
		return nil, r.attachmentMissingErr
	}
	err = r.objectUploader.PutObject(*asset.ThumbnailKey, assetThumbnail)
	if err != nil {
		return nil, err
	}
	asset.ThumbnailKey, err = r.objectUploader.GetObjectLink(*asset.ThumbnailKey)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *PsqlAssetRepository) FindAssetByID(id int) (*models.Asset, error) {
	var asset models.Asset
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE asset.assets SET 
		views=views+1
		WHERE id = $1
		RETURNING
		(id, name, description, tags, created_at, updated_at, version, upvotes, downvotes, rating, views, object_key, thumbnail_key)`,
		id,
	).Scan(&asset)
	if err != nil {
		return nil, err
	}

	asset.Key, err = r.objectUploader.GetObjectLink(*asset.Key)
	if err != nil {
		return nil, err
	}
	asset.ThumbnailKey, err = r.objectUploader.GetObjectLink(*asset.ThumbnailKey)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *PsqlAssetRepository) FindAssets(keywords []string, limit uint64, order string) (*[]models.Asset, error) {
	var assets []models.Asset

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var rows pgx.Rows
	if len(keywords) != 0 {
		query :=
			`SELECT (id, name, description, tags, created_at, updated_at, version, upvotes, downvotes, rating, views, object_key, thumbnail_key) 
				FROM
				((SELECT id, name, description, tags, created_at, updated_at, version, upvotes, downvotes, rating, views, object_key, thumbnail_key
				FROM asset.assets
				WHERE asset_ts @@ to_tsquery_multilang($1))
			UNION
				(SELECT id, name, description, tags, created_at, updated_at, version, upvotes, downvotes, rating, views, object_key, thumbnail_key 
				FROM asset.assets
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
		query := `SELECT 
			(id, name, description, tags, created_at, updated_at, version, upvotes, downvotes, rating, views, object_key, thumbnail_key) 
			FROM asset.assets`

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
		rows, err = conn.Query(context.Background(), query)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()
	for rows.Next() {
		var asset models.Asset
		err = rows.Scan(&asset)
		if err != nil {
			return nil, err
		}
		asset.Key, _ = r.objectUploader.GetObjectLink(*asset.Key)
		asset.ThumbnailKey, _ = r.objectUploader.GetObjectLink(*asset.ThumbnailKey)
		assets = append(assets, asset)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		return nil, r.NotFoundErr()
	}

	return &assets, nil
}

func (r *PsqlAssetRepository) UpdateAsset(id int, asset models.Asset, assetFile, assetThumbnail io.Reader) (*models.Asset, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(), `SELECT id FROM asset.assets WHERE id = $1 AND version = $2`, id, *asset.Version).Scan(&id)
	if err != nil {
		return nil, r.ConflictErr()
	}

	err = conn.QueryRow(context.Background(),
		`UPDATE asset.assets SET 
		name=COALESCE($1, name), description=COALESCE($2, description), tags=COALESCE($3, tags), updated_at=NOW(), upvotes=COALESCE($4, upvotes), downvotes=COALESCE($5, downvotes)
			WHERE id = $6
		RETURNING
			(id, name, description, tags, created_at, updated_at, version, upvotes, downvotes, rating, views, object_key, thumbnail_key)`,
		asset.Name, asset.Description, asset.Tags, asset.Upvotes, asset.Downvotes,
		id,
	).Scan(&asset)
	if err != nil {
		return nil, err
	}

	if assetFile != nil {
		err = r.objectUploader.PutObject(*asset.Key, assetFile)
		if err != nil {
			return nil, err
		}
	}

	if assetThumbnail != nil {
		err = r.objectUploader.PutObject(*asset.ThumbnailKey, assetThumbnail)
		if err != nil {
			return nil, err
		}
	}

	asset.Key, _ = r.objectUploader.GetObjectLink(*asset.Key)
	asset.ThumbnailKey, _ = r.objectUploader.GetObjectLink(*asset.ThumbnailKey)

	return &asset, nil
}

func (r *PsqlAssetRepository) DeleteAsset(id int) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	err = r.objectUploader.DeleteObject(fmt.Sprintf("asset-%v", id))
	if err != nil {
		return err
	}
	err = r.objectUploader.DeleteObject(fmt.Sprintf("asset-thumbnail-%v", id))
	if err != nil {
		return err
	}

	ct, err := conn.Exec(context.Background(), `DELETE FROM asset.assets WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	return nil
}
