package psqlRepository

import (
	"context"

	"gamehangar/internal/domain/models"
)

type PsqlAssetRepository struct {
	databaseClient psqlDatabaseClient
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlAssetRepository(dbClient psqlDatabaseClient) *PsqlAssetRepository {
	return &PsqlAssetRepository{
		databaseClient: dbClient,
	}
}

func (r *PsqlAssetRepository) NotFoundErr() error { return r.databaseClient.ErrNoRows() }

func (r *PsqlAssetRepository) CreateAsset(asset models.Asset) (*models.Asset, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO asset.assets
		(name, description, link, "created_at") 
		VALUES
		($1, $2, $3, $4)
		RETURNING
		(id, name, description, link, "created_at")`,
		asset.Name, asset.Description, asset.Link, asset.CreatedAt,
	).Scan(&asset)
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
		`SELECT (id, name, description, link, "created_at") 
		FROM asset.assets WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *PsqlAssetRepository) FindAssets() (*[]models.Asset, error) {
	var assets []models.Asset

	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT (id, name, description, link, "created_at") 
		FROM asset.assets`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var asset models.Asset
		err = rows.Scan(&asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &assets, nil
}

func (r *PsqlAssetRepository) UpdateAsset(id int, asset models.Asset) (*models.Asset, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE asset.assets SET 
		name=COALESCE($1, name), description=COALESCE($2, description), link=COALESCE($3, link), 
		"created_at"=COALESCE($4, "created_at")
			WHERE id = $5
		RETURNING
			(id, name, description, link, "created_at")`,
		asset.Name, asset.Description, asset.Link, asset.CreatedAt,
		id,
	).Scan(&asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *PsqlAssetRepository) DeleteAsset(id int) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM asset.assets WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		if err != nil {
			return err
		}
		return r.databaseClient.ErrNoRows()
	}
	return nil
}
