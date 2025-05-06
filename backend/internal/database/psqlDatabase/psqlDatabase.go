package psqlDatabase

import (
	"context"
	"embed"
	"gamehangar/internal/config/psqlDatabseConfig"
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
	migrator   PsqlDatabaseMigrator
	ConnPool   *pgxpool.Pool
}

type PsqlDatabaseMigrator interface {
	NewMigrator(migrationFiles any, migrationRoot, versionTable string) (any, error)
	MigrateDatabase(conn *pgx.Conn, expectedVersion int) error
}

func (p PsqlDatabase) NewDatabaseClient(connstring string, migrator PsqlDatabaseMigrator, config *psqlDatabseConfig.PsqlDatabaseConfig) (*PsqlDatabaseClient, error) {
	var err error
	var dbClient = &PsqlDatabaseClient{
		config:     config,
		connstring: connstring,
	}

	m, err := migrator.NewMigrator(MigrationFiles, config.MigrationsRoot, config.VersionTable)
	if err != nil {
		return nil, err
	}
	dbClient.migrator = m.(PsqlDatabaseMigrator)

	err = p.setup(dbClient)
	if err != nil {
		return nil, err
	}
	return dbClient, nil
}

func (p PsqlDatabase) setup(c *PsqlDatabaseClient) error {
	var err error

	c.ConnPool, err = pgxpool.New(context.Background(), c.connstring)
	if err != nil {
		return err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err = c.ConnPool.Ping(ctx)
	if err != nil {
		return err
	}

	// Migrate database
	if c.config.MigrateDatabse {
		err := c.autoMigrate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *PsqlDatabaseClient) autoMigrate() error {
	conn, err := c.ConnPool.Acquire(context.Background())
	defer conn.Release()
	if err != nil {
		return err
	}

	return c.migrator.MigrateDatabase(conn.Conn(), c.config.ExpectedVersion)
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
