package psqlDatabase

import (
	"gamehangar/internal/config/psqlDatabseConfig"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// Integration test: Requires a working PostgreSQL client and a .env file!
func TestNewDatabaseClient(t *testing.T) {
	wd, _ := os.Getwd()
	err := godotenv.Load(wd + "/../../../.env")
	if err != nil {
		t.Errorf("Error loading .env file: %v", err)
	}
	databaseConfig, err := psqlDatabseConfig.PsqlConfig{}.NewConfig(
		MigrationFiles, os.Getenv("PSQL_MIGRATE_ROOT_DIR"),
	)
	if err != nil {
		t.Errorf("Error loading PSQL database Config: %v", err)
	}
	_, err = PsqlDatabase{}.NewDatabaseClient(
		os.Getenv("PSQL_CONNSTRING"), databaseConfig,
	)
	if err != nil {
		t.Errorf("Error setting up new DatabaseClient: %v", err)
	}
}
