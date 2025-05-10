package psqlRepository

import (
	// "context"
	"errors"
	// "gamehangar/internal/config/psqlDatabseConfig"
	// "gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	// "gamehangar/pkg/ternMigrate"
	// "os"
	"testing"

	"github.com/google/uuid"
	// "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	// testDBClient     *psqlDatabase.PsqlDatabaseClient

	roleID          uuid.UUID
	roleName        string      = "Test Role"
	roleNameUpdated string      = "Test UPDATE Role"
	roleVersion     int         = 1
	role            models.Role = models.Role{Name: &roleName}
	roleUpdated     models.Role = models.Role{Name: &roleNameUpdated, Version: &roleVersion}

	userID           uuid.UUID
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
	if independent {
		ResetDB()
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
	_, err := r.FindRoleByID(uuid.New()) // (near) impossible to match
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
	_, err := r.FindUserByID(uuid.New()) // (near) impossible to match
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
	_, err := r.FindSessionByID(uuid.New()) // (near) impossible to match
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
