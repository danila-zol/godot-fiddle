package psqlRepository

import (
	"context"
	"errors"

	"gamehangar/internal/domain/models"
)

type PsqlAssetRepository struct {
	databaseClient psqlDatabaseClient
	notFoundErr    error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlAssetRepository(dbClient psqlDatabaseClient) (*PsqlAssetRepository, error) {
	return &PsqlAssetRepository{
		databaseClient: dbClient,
		notFoundErr:    errors.New("Not Found"),
	}, nil
}

func (par *PsqlAssetRepository) CreateAsset(asset models.Asset) (*models.Asset, error) {
	conn, err := par.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO asset.assets
		(id, name, description, link, createdAt) 
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING
		(id, name, description, link, createdAt)`,
		asset.ID, asset.Name, asset.Description, asset.Link, asset.CreatedAt,
	).Scan(&asset)

	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (par *PsqlAssetRepository) FindAssetByID(id string) (*models.Asset, error) {
	var asset models.Asset
	conn, err := par.databaseClient.AcquireConn()
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

func (par *PsqlAssetRepository) FindAssets() (*[]models.Asset, error) {
	var assets []models.Asset

	conn, err := par.databaseClient.AcquireConn()
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

func (par *PsqlAssetRepository) UpdateAsset(id string, asset models.Asset) (*models.Asset, error) {
	conn, err := par.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`UPDATE asset.assets SET 
		name=$1, description=$2, link=$3, createdAt=$4
		WHERE id = $5
		RETURNING
			(id, name, description, link, createdAt)`,
		asset.Name, asset.Description, asset.Link, asset.CreatedAt,
		id,
	).Scan(&asset)
	// TODO: How to handle 404 error here?
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (par *PsqlAssetRepository) DeleteAsset(id string) error {
	conn, err := par.databaseClient.AcquireConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM asset.assets WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return par.notFoundErr
	}
	if err != nil {
		return err
	}
	return nil
}
