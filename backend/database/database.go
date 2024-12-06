package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

//go:embed data/*.sql
var migrationFiles embed.FS

func migrateDatabase(conn *pgx.Conn, schemaVersion string) {
	migrator, err := migrate.NewMigrator(context.Background(), conn, schemaVersion)
	if err != nil {
		log.Fatalf("Unable to create a migrator: %v\n", err)
	}

	migrationRoot, err := fs.Sub(migrationFiles, "data")
	err = migrator.LoadMigrations(migrationRoot)
	if err != nil {
		log.Fatalf("Unable to load migrations: %v\n", err)
	}

	err = migrator.Migrate(context.Background())
	if err != nil {
		log.Fatalf("Unable to migrate: %v\n", err)
	}

	ver, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		log.Fatalf("Unable to get current schema version: %v\n", err)
	}

	log.Printf("Migration done. Current schema version: %v\n", ver)
}

func SetupDB(dsn string) {
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v\n", err)
	}

	autoMigrate, err := strconv.ParseBool(os.Getenv("MIGRATE_DB"))
	if err == nil {
		if autoMigrate == true {
			migrateDatabase(conn.Conn(), os.Getenv("MIGRATE_SCHEMA"))
		}
	} else {
		log.Println("WARNING! Missing env variable DB_MIGRATE. Assuming false")
	}
}
