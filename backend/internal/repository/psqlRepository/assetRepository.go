package psqlRepository

import (
	"context"

	"gamehangar/internal/database/psql"
	"gamehangar/internal/domain/models"
)

type PsqlAssetRepository struct {
	databaseClient *psql.PsqlDatabaseClient
}

func (par *PsqlAssetRepository) CreateAsset(asset models.Asset) error {
	conn, err := par.databaseClient.ConnPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO asset.assets
		(id, name, description, link, createdAt) 
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING
		(id, name, description, link, createdAt)`,
		asset.ID, asset.Name, asset.Description, asset.Link, asset.CreatedAt,
	)

	err = row.Scan(&asset)
	if err != nil {
		return err
	}
	return nil
}
