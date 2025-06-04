package services

import (
	"context"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	"gamehangar/internal/enforcer/psqlCasbinClient"
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
	testUploader *ObjectUploader
	testDBClient *psqlDatabase.PsqlDatabaseClient
	testEnforcer *psqlCasbinClient.CasbinClient

	userAuthorizer *UserAuthorizer

	role           string = "admin"
	userID         uuid.UUID
	userRepository *psqlRepository.PsqlUserRepository

	testUsername string = "danila-zol"
	testEmail    string = "tset@email.com"
	testPassword string = "aVeryStrongPassword123"
)

func init() {
	var err error
	wd, _ := os.Getwd()
	if independent {
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

		CREATE TABLE "user".users (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"username" VARCHAR(255) NOT NULL UNIQUE,
		"display_name" VARCHAR(255),
		"email" VARCHAR(255) NOT NULL UNIQUE,
		"password" VARCHAR(255) NOT NULL,
		"verified" BOOLEAN NOT NULL DEFAULT false,
		"role" VARCHAR(255) NOT NULL,
		"created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		"karma" INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE "user".sessions (
		"id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
		"user_id" UUID NOT NULL REFERENCES "user".users (id) ON DELETE CASCADE
		);
		`)
	if err != nil {
		panic("Error resetting user schema" + err.Error())
	}
	err = c.QueryRow(
		context.Background(),
		`INSERT INTO "user".users
		(username, display_name, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING (id)`,
		testUsername, "Mike", testEmail, "TestPassword", role).Scan(&userID)
	if err != nil {
		panic("Error resetting user schema" + err.Error())
	}
	testEnforcer, err = psqlCasbinClient.CasbinConfig{}.NewCasbinClient(
		os.Getenv("PSQL_CONNSTRING"),
		wd+"/../enforcer/psqlCasbinClient/rbac_model.conf",
	)
	if err != nil {
		panic("Error creating Casbin enforcer: " + err.Error())
	}
	userRepository = psqlRepository.NewPsqlUserRepository(testDBClient, testEnforcer, testS3Client)
	userAuthorizer = NewUserAuthorizer(userRepository, testEnforcer)
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

	_, err = userRepository.UpdateUser(userID, models.User{Password: hash}, nil)
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
	err = r.DeleteRole(role)
	if err != nil {
		panic(err)
	}
}
