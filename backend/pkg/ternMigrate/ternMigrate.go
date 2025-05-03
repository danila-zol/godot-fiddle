package ternMigrate

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

// A simple tern Migrator wrapper
type Migrator struct {
	migrationFiles embed.FS
	migrationRoot  string
	versionTable   string
}

func (m Migrator) NewMigrator(migrationFiles any, migrationRoot, versionTable string) (any, error) {
	mf, ok := migrationFiles.(embed.FS)
	if !ok {
		return nil, errors.New("Failed to migrate PSQL Database: Invalid migration files!")
	}
	return &Migrator{
		migrationFiles: mf,
		migrationRoot:  migrationRoot,
		versionTable:   versionTable,
	}, nil
}

func (m Migrator) MigrateDatabase(conn *pgx.Conn, expected int) error {
	migrator, err := migrate.NewMigrator(context.Background(), conn, m.versionTable)
	if err != nil {
		return err
	}

	migrationRoot, err := fs.Sub(m.migrationFiles, m.migrationRoot)
	if err != nil {
		return err
	}
	err = migrator.LoadMigrations(migrationRoot)
	if err != nil {
		return err
	}

	now, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		return err
	}
	if now != int32(expected) {
		log.Printf("Current version: %v\nExpected version: %v\nMigrating...\n", now, expected)
		err = migrator.MigrateTo(context.Background(), int32(expected))
		if err != nil {
			return err
		}

	}

	ver, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		return err
	}

	log.Printf("Migration done. Current schema version: %v\n", ver)
	return nil
}
