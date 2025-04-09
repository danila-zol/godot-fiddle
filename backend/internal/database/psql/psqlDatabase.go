package psql

import (
	"context"
	"embed"
	"errors"
	"gamehangar/internal/config"
	"gamehangar/pkg/ternMigrate"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var MigrationFiles embed.FS

type PsqlDatabaseClient struct {
	config     *config.PsqlDatabaseConfig
	connstring string
	ConnPool   *pgxpool.Pool // nil until Setup() is called
}

func NewDatabaseClient(connstring string, config *config.PsqlDatabaseConfig) *PsqlDatabaseClient {
	var dbClient = &PsqlDatabaseClient{
		config:     config,
		connstring: connstring,
	}

	return dbClient
}

func (pdc *PsqlDatabaseClient) Setup() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), pdc.connstring)
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	// Migrate database
	if pdc.config.MigrateDatabse {
		err := pdc.autoMigrate()
		if err != nil {
			return nil, err
		}
	}

	return dbpool, nil
}

func (pdc *PsqlDatabaseClient) autoMigrate() error {
	migrationFiles, ok := pdc.config.Migrations.(embed.FS)
	if !ok {
		return errors.New("Failed to migrate PSQL Database: Invalid migration files!")
	}
	migrator := &ternMigrate.Migrator{
		MigrationFiles: migrationFiles,
		MigrationDir:   pdc.config.MigrationsDir,
		VersionTable:   pdc.config.VersionTable,
	}
	// TODO: Fix logic! ConnPool is nil until Setup() is complete!
	conn, err := pdc.ConnPool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		return err
	}
	migrator.MigrateDatabase(conn.Conn(), int32(pdc.config.ExpectedVersion))
	return nil
}
