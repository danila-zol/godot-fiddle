package config

import (
	"os"
	"strconv"
)

type PsqlDatabaseConfig struct {
	MigrateDatabse  bool
	Migrations      any
	MigrationsDir   string // TODO: Relative to where?
	ExpectedVersion int
	VersionTable    string
}

func NewConfig(migrations any, migrationsDir string) (*PsqlDatabaseConfig, error) {
	migrateDatabase, err := strconv.ParseBool(os.Getenv("PSQL_MIGRATE_DATABASE"))
	if err != nil {
		return nil, err
	}

	versionTable := os.Getenv(
		"PSQL_MIGRATE_VERSION_TABLE")
	expectedVersion, err := strconv.ParseInt(os.Getenv("PSQL_MIGRATE_EXPECTED_VERSION"), 10, 32)
	if err != nil {
		return nil, err
	}

	return &PsqlDatabaseConfig{
		MigrateDatabse:  migrateDatabase,
		Migrations:      migrations,
		MigrationsDir:   migrationsDir,
		ExpectedVersion: int(expectedVersion),
		VersionTable:    versionTable,
	}, nil
}
