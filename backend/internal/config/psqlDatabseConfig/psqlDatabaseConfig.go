package psqlDatabseConfig

import (
	"os"
	"strconv"
)

type PsqlConfig struct{} // PostgreSQL-related configs and mehtods

type PsqlDatabaseConfig struct {
	MigrateDatabse  bool
	Migrations      any
	MigrationsRoot  string
	ExpectedVersion int
	VersionTable    string
}

func (p PsqlConfig) NewConfig(migrations any, migrationsRoot string) (*PsqlDatabaseConfig, error) {
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
		MigrationsRoot:  migrationsRoot,
		ExpectedVersion: int(expectedVersion),
		VersionTable:    versionTable,
	}, nil
}
