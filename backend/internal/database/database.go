package database

import "embed"

//go:embed migrations/*.sql
var migrationFiles embed.FS

type DatabaseClient struct {
	Migrations embed.FS
	Connstring string
}

func NewDatabaseClient(connstring string) (*DatabaseClient, error) {
	var dbClient = &DatabaseClient{
		Migrations: migrationFiles,
		Connstring: connstring,
	}

	return dbClient, nil
}
