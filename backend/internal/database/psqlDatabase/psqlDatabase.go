package psqlDatabase

import (
	"context"
	"embed"
	"errors"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/pkg/ternMigrate"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var MigrationFiles embed.FS

type PsqlDatabase struct{} // PostgreSQL-related database clients and methods

type PsqlDatabaseClient struct {
	config     *psqlDatabseConfig.PsqlDatabaseConfig
	connstring string
	ConnPool   *pgxpool.Pool
}

func (p PsqlDatabase) NewDatabaseClient(connstring string, config *psqlDatabseConfig.PsqlDatabaseConfig) (*PsqlDatabaseClient, error) {
	var dbClient = &PsqlDatabaseClient{
		config:     config,
		connstring: connstring,
	}

	err := p.setup(dbClient)
	if err != nil {
		return nil, err
	}
	return dbClient, nil
}

func (p PsqlDatabase) setup(pdc *PsqlDatabaseClient) error {
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

func (c *PsqlDatabaseClient) autoMigrate() error {
	migrationFiles, ok := c.config.Migrations.(embed.FS)
	if !ok {
		return errors.New("Failed to migrate PSQL Database: Invalid migration files!")
	}
	migrator := &ternMigrate.Migrator{
		MigrationFiles: migrationFiles,
		MigrationRoot:  c.config.MigrationsRoot,
		VersionTable:   c.config.VersionTable,
	}
	conn, err := c.ConnPool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		return err
	}
	migrator.MigrateDatabase(conn.Conn(), int32(c.config.ExpectedVersion))
	return nil
}

func (c *PsqlDatabaseClient) AcquireConn() (*pgxpool.Conn, error) {
	conn, err := c.ConnPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *PsqlDatabaseClient) ErrNoRows() error {
	return pgx.ErrNoRows
}
