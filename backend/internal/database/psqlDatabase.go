package database

import (
	"context"
	"embed"
	"errors"
	"gamehangar/internal/config"
	"gamehangar/pkg/ternMigrate"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed psqlMigrations/*.sql
var MigrationFiles embed.FS

type Psql struct{} // PostgreSQL-related database clients and methods

type PsqlDatabaseClient struct {
	config     *config.PsqlDatabaseConfig
	connstring string
	ConnPool   *pgxpool.Pool // nil until Setup() is called
}

func (p Psql) NewDatabaseClient(connstring string, config *config.PsqlDatabaseConfig) *PsqlDatabaseClient {
	var dbClient = &PsqlDatabaseClient{
		config:     config,
		connstring: connstring,
	}

	return dbClient
}

func (p Psql) Setup(pdc *PsqlDatabaseClient) error {
	var err error

	pdc.ConnPool, err = pgxpool.New(context.Background(), pdc.connstring)
	if err != nil {
		return err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err = pdc.ConnPool.Ping(ctx)
	if err != nil {
		return err
	}

	// Migrate database
	if pdc.config.MigrateDatabse {
		err := pdc.autoMigrate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (pdc *PsqlDatabaseClient) autoMigrate() error {
	migrationFiles, ok := pdc.config.Migrations.(embed.FS)
	if !ok {
		return errors.New("Failed to migrate PSQL Database: Invalid migration files!")
	}
	migrator := &ternMigrate.Migrator{
		MigrationFiles: migrationFiles,
		MigrationRoot:  pdc.config.MigrationsRoot,
		VersionTable:   pdc.config.VersionTable,
	}
	conn, err := pdc.ConnPool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		return err
	}
	migrator.MigrateDatabase(conn.Conn(), int32(pdc.config.ExpectedVersion))
	return nil
}
