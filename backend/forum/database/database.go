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
var dbpool *pgxpool.Pool

func migrateDatabase(conn *pgx.Conn, versionColumn string, expected int32) {
	migrator, err := migrate.NewMigrator(context.Background(), conn, versionColumn)
	if err != nil {
		log.Fatalf("Unable to create a migrator: %v\n", err)
	}

	migrationRoot, err := fs.Sub(migrationFiles, "data")
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

func SetupDB(dsn string) {
	var err error
	dbpool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	conn, err := dbpool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v\n", err)
	}

	autoMigrate, err := strconv.ParseBool(os.Getenv("MIGRATE_DB"))
	if err == nil {
		if autoMigrate == true {
			expected, err := strconv.ParseInt(os.Getenv("EXPECTED_VERSION"), 10, 32)
			if err != nil {
				log.Fatalf("Unable to get database expected version: %v\n", err)
			}
			migrateDatabase(
				conn.Conn(),
				os.Getenv("VERSION_COLUMN"),
				(int32(expected)),
			)
		}
	} else {
		log.Printf("WARNING! Missing env variable DB_MIGRATE: %v\nAssuming false\n", err)
	}
}
