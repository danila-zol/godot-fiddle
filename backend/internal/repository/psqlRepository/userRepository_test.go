package psqlRepository

import (
	"context"
	"errors"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	// testDBClient     *psqlDatabase.PsqlDatabaseClient

	roleID          string
	roleName        string      = "Test Role"
	roleNameUpdated string      = "Test UPDATE Role"
	roleVersion     int         = 1
	role            models.Role = models.Role{Name: &roleName}
	roleUpdated     models.Role = models.Role{Name: &roleNameUpdated, Version: &roleVersion}

	userID           string
	userName         string      = "Test User"
	userNameUpdated  string      = "Test UPDATE User"
	userDisplayName  string      = "A user for integration testing for PSQL Repo (yes, a display name)"
	userEmail        string      = "user@example.com"
	userPassword     string      = "verySecurePassword"
	userKarmaUpdated int         = 9
	user             models.User = models.User{Username: &userName, Email: &userEmail, Password: &userPassword}
	userUpdated      models.User = models.User{Username: &userNameUpdated, DisplayName: &userDisplayName, Karma: &userKarmaUpdated}

	session models.Session
)

func init() {
	wd, _ := os.Getwd()
	err := godotenv.Load(wd + "/../../../.env")
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
		os.Getenv("PSQL_CONNSTRING"), databaseConfig,
	)
	if err != nil {
		panic("Error setting up new DatabaseClient")
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
}

func TestCreateRole(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	resultRole, err := r.CreateRole(role)
	assert.NoError(t, err)
	role = *resultRole
}

func TestFindRoleByID(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	_, err := r.FindRoleByID(*role.ID)
	assert.NoError(t, err)
}

func TestFindRoleByIDNoRows(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	_, err := r.FindRoleByID("c6c00b97-0264-4c4a-b01e-e2f2a3f02572") // (near) impossible to match
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestUpdateRole(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	resultRole, err := r.UpdateRole(*role.ID, roleUpdated)
	assert.NoError(t, err)

	modifiedRole := role
	modifiedRole.Name = &roleNameUpdated
	newVersion := *roleUpdated.Version + 1
	modifiedRole.Version = &newVersion

	assert.Equal(t, modifiedRole, *resultRole)
}

func TestUpdateRoleMultiple(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}

	modifiedRole := role // Manual update
	modifiedRole.ID = role.ID

	for i := 2; i < 6; i++ {
		newName := *roleUpdated.Name + " New"
		roleUpdated.Name = &newName
		newVersion := i
		roleUpdated.Version = &newVersion

		resultRole, err := r.UpdateRole(*role.ID, roleUpdated)
		assert.NoError(t, err)

		newerVersion := i + 1
		modifiedRole.Version = &newerVersion
		modifiedRole.Name = &newName

		assert.Equal(t, modifiedRole, *resultRole)
	}
}

func TestUpdateRoleConflict(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.UpdateRole(*role.ID, roleUpdated)
	if assert.Error(t, err) {
		assert.Equal(t, r.conflictErr, err)
	}
}

func TestDeleteRole(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	err := r.DeleteRole(*role.ID)
	assert.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	role, err := r.CreateRole(role)
	roleID = *role.ID
	user.RoleID = role.ID

	resultUser, err := r.CreateUser(user)
	assert.NoError(t, err)
	user = *resultUser
}

func TestFindUserByID(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindUserByID(*user.ID)
	assert.NoError(t, err)
}

func TestFindUserByEmail(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindUserByEmail(*user.Email)
	assert.NoError(t, err)
}

func TestFindUserByUsername(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, conflictErr: errors.New("Record conflict!")}
	_, err := r.FindUserByUsername(*user.Username)
	assert.NoError(t, err)
}

func TestFindUserByIDNoRows(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	_, err := r.FindUserByID("c6c00b97-0264-4c4a-b01e-e2f2a3f02572") // (near) impossible to match
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestUpdateUser(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	resultUser, err := r.UpdateUser(*user.ID, userUpdated)
	assert.NoError(t, err)

	modifiedUser := user
	modifiedUser.Username = &userNameUpdated
	modifiedUser.DisplayName = &userDisplayName
	modifiedUser.Karma = &userKarmaUpdated

	assert.Equal(t, modifiedUser, *resultUser)
}

func TestDeleteUser(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	err := r.DeleteUser(*user.ID)
	assert.NoError(t, err)
}

func TestCreateSession(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	user, err := r.CreateUser(user)
	userID = *user.ID
	session.UserID = user.ID

	resultSession, err := r.CreateSession(session)
	assert.NoError(t, err)
	session = *resultSession
}

func TestFindSessionByID(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	_, err := r.FindSessionByID(*session.ID)
	assert.NoError(t, err)
}

func TestFindSessionByIDNoRows(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	_, err := r.FindSessionByID("c6c00b97-0264-4c4a-b01e-e2f2a3f02572") // (near) impossible to match
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestDeleteSession(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient}
	err := r.DeleteSession(*session.ID)
	if assert.NoError(t, err) {
		teardownUser(&r)
	}
}

func teardownUser(r *PsqlUserRepository) {
	err := r.DeleteUser(userID)
	if err != nil {
		panic(err)
	}
	err = r.DeleteRole(roleID)
	if err != nil {
		panic(err)
	}
}
