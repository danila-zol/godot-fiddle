package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// TODO: MapNameToUser + assert one unique user per username!
// TODO: A new asset creates a new topic on forums — a different service!

func CreateAsset(asset Asset) (*Asset, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO asset.assets
			(id, name, description, link, created_at) 
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			(id, name, description, link, created_at)`,
		asset.ID, asset.Name, asset.Description, asset.Link, asset.Created_at,
	)

	err = row.Scan(&asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func FindFirstAsset(id string) (*Asset, error) {
	var asset Asset
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT * FROM asset.assets WHERE id = $1 LIMIT 1`,
		id,
	).Scan(&asset.ID, &asset.Name, &asset.Description, &asset.Link, &asset.Created_at)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func FindAssets() (*[]Asset, error) {
	var assets []Asset

	conn, err := dbpool.Acquire(context.Background())
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
		var asset Asset
		err = rows.Scan(&asset.ID, &asset.Name, &asset.Description, &asset.Link, &asset.Created_at)
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

func UpdateAsset(asset Asset) (*Asset, error) {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(),
		`UPDATE asset.assets SET 
		name=$1, description=$2, link=$3, created_at=$4
		WHERE id = $5`,
		asset.Name, asset.Description, asset.Link, asset.Created_at,
		asset.ID,
	)
	if ct.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func DeleteAsset(id string) error {
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	ct, err := conn.Exec(context.Background(), `DELETE FROM asset.assets WHERE id=$1`, id)
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if err != nil {
		return err
	}
	return nil
}
