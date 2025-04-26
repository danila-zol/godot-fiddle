package psqlRepository

import (
	"context"
	"errors"

	"gamehangar/internal/domain/models"
)

type PsqlAssetRepository struct {
	databaseClient psqlDatabaseClient
	conflictErr    error
}

// Requires PsqlDatabaseClient since it implements PostgeSQL-specific query logic
func NewPsqlAssetRepository(dbClient psqlDatabaseClient) *PsqlAssetRepository {
	return &PsqlAssetRepository{
		databaseClient: dbClient,
		conflictErr:    errors.New("Record conflict!"),
	}
}

func (r *PsqlAssetRepository) NotFoundErr() error { return r.databaseClient.ErrNoRows() }

// Returns "Record conflict!" to specify conflicting record versions on update
func (r *PsqlAssetRepository) ConflictErr() error { return r.conflictErr }

func (r *PsqlAssetRepository) CreateAsset(asset models.Asset) (*models.Asset, error) {
	conn, err := r.databaseClient.AcquireConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`INSERT INTO asset.assets
		(name, description, link) 
		VALUES
		($1, $2, $3)
		RETURNING
		(id, name, description, link, created_at, version)`,
		asset.Name, asset.Description, asset.Link,
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
		`SELECT (id, name, description, link, created_at, version) 
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
		`SELECT (id, name, description, link, created_at, version) 
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

	ct, err := conn.Exec(context.Background(), `SELECT (id) FROM asset.assets WHERE id = $1 AND version = $2`, id, *asset.Version)
	if err != nil || ct.RowsAffected() == 0 {
		return nil, r.ConflictErr()
	}

	err = conn.QueryRow(context.Background(), // TODO: Create a sequence to increment on update
		`UPDATE asset.assets SET 
		name=COALESCE($1, name), description=COALESCE($2, description), link=COALESCE($3, link), version=version+1
			WHERE id = $4
		RETURNING
			(id, name, description, link, created_at, version)`,
		asset.Name, asset.Description, asset.Link,
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
