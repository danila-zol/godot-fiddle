package ternMigrate

import (
	"context"
	"embed"
	"io/fs"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

// A simple tern Migrator wrapper
type Migrator struct {
	MigrationFiles embed.FS
	MigrationRoot  string
	VersionTable   string
}

func (m Migrator) MigrateDatabase(conn *pgx.Conn, expected int32) {
	migrator, err := migrate.NewMigrator(context.Background(), conn, m.VersionTable)
	if err != nil {
		log.Fatalf("Unable to create a migrator: %v\n", err)
	}

	migrationRoot, err := fs.Sub(m.MigrationFiles, m.MigrationRoot)
	if err != nil {
		log.Fatalf("Unable to load migrations: %v\n", err)
	}
	err = migrator.LoadMigrations(migrationRoot)
	if err != nil {
		log.Fatalf("Unable to load migrations: %v\n", err)
	}

	now, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		log.Fatalf("Unable to get current schema version: %v\n", err)
	}
	if now != expected {
		log.Printf("Current version: %v\nExpected version: %v\nMigrating...\n", now, expected)
		err = migrator.MigrateTo(context.Background(), expected)
		if err != nil {
			log.Fatalf("Unable to migrate: %v\n", err)
		}

	}

	ver, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		log.Fatalf("Unable to get current schema version: %v\n", err)
	}

	log.Printf("Migration done. Current schema version: %v\n", ver)
}
