package psqlRepository

import (
	// "context"
	// "gamehangar/internal/config/psqlDatabseConfig"
	// "gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/domain/models"
	"gamehangar/internal/enforcer/psqlCasbinClient"
	// "gamehangar/pkg/ternMigrate"
	// "os"
	"testing"

	"github.com/google/uuid"
	// "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	// testDBClient     *psqlDatabase.PsqlDatabaseClient
	// testS3Client *MockS3
	testEnforcer *psqlCasbinClient.CasbinClient

	role string = "Sharif"

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
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	err := r.CreateRole(role)
	assert.NoError(t, err)
}
func TestDeleteRole(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	err := r.DeleteRole(role)
	assert.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	err := r.CreateRole(role)
	user.Role = &role

	resultUser, err := r.CreateUser(user, nil)
	assert.NoError(t, err)
	user = *resultUser
}

func TestFindUserByID(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	_, err := r.FindUserByID(*user.ID)
	assert.NoError(t, err)
}

func TestFindUserByEmail(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	_, err := r.FindUserByEmail(*user.Email)
	assert.NoError(t, err)
}

func TestFindUserByUsername(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	_, err := r.FindUserByUsername(*user.Username)
	assert.NoError(t, err)
}

func TestFindUserByIDNoRows(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer}
	_, err := r.FindUserByID(uuid.New()) // (near) impossible to match
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestUpdateUser(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	resultUser, err := r.UpdateUser(*user.ID, userUpdated, nil)
	assert.NoError(t, err)

	modifiedUser := user
	modifiedUser.Username = &userNameUpdated
	modifiedUser.DisplayName = &userDisplayName
	modifiedUser.Karma = &userKarmaUpdated

	l := "https://example.com"
	modifiedUser.ProfilePic = &l

	assert.Equal(t, modifiedUser, *resultUser)
}

func TestDeleteUser(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	err := r.DeleteUser(*user.ID)
	assert.NoError(t, err)
}

func TestCreateSession(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	user, err := r.CreateUser(user, nil)
	userID = *user.ID
	session.UserID = user.ID

	resultSession, err := r.CreateSession(session)
	assert.NoError(t, err)
	session = *resultSession
}

func TestFindSessionByID(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	_, err := r.FindSessionByID(*session.ID)
	assert.NoError(t, err)
}

func TestFindSessionByIDNoRows(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
	_, err := r.FindSessionByID(uuid.New()) // (near) impossible to match
	if assert.Error(t, err) {
		assert.Equal(t, r.NotFoundErr(), err)
	}
}

func TestDeleteSession(t *testing.T) {
	r := PsqlUserRepository{databaseClient: testDBClient, enforcer: testEnforcer, objectUploader: testS3Client}
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
	err = r.DeleteRole(role)
	if err != nil {
		panic(err)
	}
}
