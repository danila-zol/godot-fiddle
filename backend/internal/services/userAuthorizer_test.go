package services

import (
	"context"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	"gamehangar/internal/repository/psqlRepository"
	"gamehangar/pkg/ternMigrate"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	// independent bool = false
	testDBClient *psqlDatabase.PsqlDatabaseClient

	userAuthorizer *UserAuthorizer

	roleID         uuid.UUID
	userID         uuid.UUID
	userRepository *psqlRepository.PsqlUserRepository

	testUsername string = "danila-zol"
	testEmail    string = "tset@email.com"
	testPassword string = "aVeryStrongPassword123"
)

func init() {
	var err error
	if independent {
		wd, _ := os.Getwd()
		err := godotenv.Load(wd + "/../../.env")
		if err != nil {
			panic("Error loading .env file:" + err.Error() + ": " + wd)
		}
		databaseConfig, err := psqlDatabseConfig.PsqlConfig{}.NewConfig(
			psqlDatabase.MigrationFiles, os.Getenv("PSQL_MIGRATE_ROOT_DIR"),
		)
		if err != nil {
			panic("Error loading PSQL database Config")
		}
		testDBClient, err = psqlDatabase.PsqlDatabase{}.NewDatabaseClient(
			os.Getenv("PSQL_CONNSTRING"), ternMigrate.Migrator{}, databaseConfig,
		)
		if err != nil {
			panic("Error setting up new DatabaseClient")
		}
	}
	c, _ := testDBClient.AcquireConn() // WARNING! Integration tests DROP TABLEs
	_, err = c.Exec(context.Background(), `
		DROP TRIGGER IF EXISTS increment_role_version_on_update ON "user".roles; 
		DROP SCHEMA IF EXISTS "user" CASCADE;

		CREATE SCHEMA IF NOT EXISTS "user";

		CREATE TABLE "user".roles (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"name" VARCHAR(255) NOT NULL,
		"version" INTEGER NOT NULL DEFAULT 1
		-- "permissions" VARCHAR(64)[]
		);

		CREATE TABLE "user".users (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"username" VARCHAR(255) NOT NULL UNIQUE,
		"display_name" VARCHAR(255),
		"email" VARCHAR(255) NOT NULL UNIQUE,
		"password" VARCHAR(255) NOT NULL,
		"verified" BOOLEAN NOT NULL DEFAULT false,
		"role_id" UUID NOT NULL REFERENCES "user".roles (id) ON DELETE RESTRICT,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"karma" INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE "user".sessions (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"user_id" UUID NOT NULL REFERENCES "user".users (id) ON DELETE CASCADE
		);

		CREATE OR REPLACE FUNCTION increment_version()
		RETURNS TRIGGER AS
		$func$
		BEGIN
		NEW.version := OLD.version + 1;
		RETURN NEW;
		END;
		$func$ LANGUAGE plpgsql;

		CREATE TRIGGER increment_role_version_on_update
		BEFORE UPDATE ON "user".roles
		FOR EACH ROW
		EXECUTE FUNCTION increment_version();
		`)
	if err != nil {
		panic("Error resetting user schema" + err.Error())
	}
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO "user".roles (name) VALUES ($1) RETURNING (id)`,
		"admin").Scan(&roleID)
	if err != nil {
		panic("Error resetting user schema" + err.Error())
	}
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO "user".users
		(username, display_name, email, password, role_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING (id)`,
		testUsername, "Mike", testEmail, "TestPassword", roleID).Scan(&userID)
	if err != nil {
		panic("Error resetting user schema" + err.Error())
	}
	userRepository = psqlRepository.NewPsqlUserRepository(testDBClient)
	userAuthorizer = NewUserAuthorizer(userRepository)
}

func TestIdentifyUserUsername(t *testing.T) {
	u, err := userAuthorizer.IdentifyUser(nil, &testUsername)
	if assert.NoError(t, err) {
		assert.Equal(t, userID, *u.ID)
	}
}

func TestIdentifyUserEmail(t *testing.T) {
	u, err := userAuthorizer.IdentifyUser(&testEmail, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, userID, *u.ID)
	}
}

func TestIdentifyUserNil(t *testing.T) {
	_, err := userAuthorizer.IdentifyUser(nil, nil)
	if assert.Error(t, err) {
		assert.Equal(t, testDBClient.ErrNoRows(), err)
	}
}

func TestPasswordHashCreateCheck(t *testing.T) {
	hash, err := userAuthorizer.CreatePasswordHash(&testPassword)
	assert.NoError(t, err)

	_, err = userRepository.UpdateUser(userID, models.User{Password: hash})
	assert.NoError(t, err)

	err = userAuthorizer.CheckPassword(&testPassword, userID)
	if assert.NoError(t, err) {
		teardownAuthorizer(userRepository)
	}
}

func teardownAuthorizer(r *psqlRepository.PsqlUserRepository) {
	err := r.DeleteUser(userID)
	if err != nil {
		panic(err)
	}
	err = r.DeleteRole(roleID)
	if err != nil {
		panic(err)
	}
}
