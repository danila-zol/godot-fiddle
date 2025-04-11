package psqlRepository

import (
	"context"

	"gamehangar/internal/domain/models"
)

type PsqlAssetRepository struct {
	databaseClient psqlDatabaseClient
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlAssetRepository(dbClient psqlDatabaseClient) (*PsqlAssetRepository, error) {
	return &PsqlAssetRepository{
		databaseClient: dbClient,
	}, nil
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
		(id, name, description, link, "createdAt") 
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING
		(id, name, description, link, "createdAt")`,
		asset.ID, asset.Name, asset.Description, asset.Link, asset.CreatedAt,
	).Scan(&asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *PsqlAssetRepository) FindAssetByID(id string) (*models.Asset, error) {
	var asset models.Asset
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM asset.assets WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&asset.ID, &asset.Name, &asset.Description, &asset.Link, &asset.CreatedAt)
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

	rows, err := conn.Query(context.Background(), `SELECT * FROM asset.assets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var asset models.Asset
		err = rows.Scan(&asset.ID, &asset.Name, &asset.Description, &asset.Link, &asset.CreatedAt)
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

func (r *PsqlAssetRepository) UpdateAsset(id string, asset models.Asset) (*models.Asset, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE asset.assets SET 
		name=COALESCE($1, name), description=COALESCE($2, description), link=COALESCE($3, link), "createdAt"=COALESCE($4, "createdAt")
		WHERE id = $5
		RETURNING
			(id, name, description, link, "createdAt")`,
		asset.Name, asset.Description, asset.Link, asset.CreatedAt,
		id,
	).Scan(&asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *PsqlAssetRepository) DeleteAsset(id string) error {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM asset.assets WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return r.databaseClient.ErrNoRows()
	}
	if err != nil {
		return err
	}
	return nil
}
