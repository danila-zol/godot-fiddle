package handlers

import (
	"errors"
	"gamehangar/internal/domain/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockUserAuthorizer struct {
	repository *mockUserRepo
}

type mockUserRepo struct {
	roleData    map[string]models.Role
	sessionData map[string]models.Session
	userData    map[string]models.User
	notFoundErr error
	conflictErr error
}

var (
	// v = validator.New(validator.WithRequiredStructEnabled())
	mu = mockUserRepo{
		roleData:    make(map[string]models.Role, 1),
		sessionData: make(map[string]models.Session, 1),
		userData:    make(map[string]models.User, 1),
		notFoundErr: errors.New("Not Found"),
		conflictErr: errors.New("Record conflict!"),
	}
	au = mockUserAuthorizer{&mu}

	genericUUID string = "9c6ac0b1-b97e-4356-a6e1-dc6b52324220"

	notFoundResponse          = `{"code":404,"message":"Not Found!"}` + "\n"
	conflictResponse          = `{"code":409,"message":"Error: unable to update the record due to an edit conflict, please try again!"}` + "\n"
	verifiedResponse          = `User verified`
	loginResponse             = `Login successful`
	logoutResponse            = `Session successfully deleted!`
	passwordResetResponse     = `User password reset!`
	passwordIncorrectResponse = `{"code":401,"message":"Password incorrect!"}` + "\n"

	roleJSON               = `{"name":"Cool role"}`
	roleJSONExpected       = `{"id":"` + genericUUID + `","name":"Cool role","version":1}` + "\n"
	roleJSONUpdate         = `{"name":"Updated cool role","version":1}`
	roleJSONUpdateInvalid  = `{"name":"Updated cool role"}`
	roleJSONUpdateExpected = `{"id":"` + genericUUID + `","name":"Updated cool role","version":2}` + "\n"

	userJSON               = `{"username":"Cool user","email":"test@example.com","roleID":"` + genericUUID + `"}`
	userJSONExpected       = `{"id":"` + genericUUID + `","username":"Cool user","email":"test@example.com","verified":false,"roleID":"` + genericUUID + `"}` + "\n"
	userJSONExpectedMany   = `[{"id":"` + genericUUID + `","username":"Cool user","email":"test@example.com","verified":false,"roleID":"` + genericUUID + `"}]` + "\n"
	userJSONUpdate         = `{"username":"Updated cool user"}`
	userJSONUpdateExpected = `{"id":"` + genericUUID + `","username":"Updated cool user","email":"test@example.com","verified":false,"roleID":"` + genericUUID + `"}` + "\n"
	userJSONVerifyExpected = `{"id":"` + genericUUID + `","username":"Updated cool user","email":"test@example.com","verified":true,"roleID":"` + genericUUID + `"}` + "\n"

	userPassword      = "qwety123"
	userPasswordReset = "aVeryStrongPassword123"
	userEmail         = "test@example.com"
	userUsername      = "Cool user"
)

func (r *mockUserRepo) CreateRole(role models.Role) (*models.Role, error) {
	id := genericUUID
	role.ID = &id
	version := 1
	role.Version = &version
	r.roleData[id] = role
	resultRole, ok := r.roleData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &resultRole, nil
}
func (r *mockUserRepo) FindRoleByID(id string) (*models.Role, error) {
	role, ok := r.roleData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &role, nil
}
func (r *mockUserRepo) UpdateRole(id string, role models.Role) (*models.Role, error) {
	var resultRole models.Role
	_, ok := r.roleData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	resultRole = r.roleData[id]
	if *resultRole.Version != *role.Version {
		return nil, r.ConflictErr()
	}
	if role.Name != nil {
		resultRole.Name = role.Name
		n := *role.Version + 1
		resultRole.Version = &n
		r.roleData[id] = resultRole
	}
	resultRole = r.roleData[id]
	return &resultRole, nil
}
func (r *mockUserRepo) DeleteRole(id string) error {
	_, ok := r.roleData[id]
	if !ok {
		return r.NotFoundErr()
	}
	delete(r.roleData, id)
	return nil
}

func (r *mockUserRepo) CreateSession(session models.Session) (*models.Session, error) {
	id := genericUUID
	session.ID = &id
	r.sessionData[id] = session
	resultsession, ok := r.sessionData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &resultsession, nil
}
func (r *mockUserRepo) FindSessionByID(id string) (*models.Session, error) {
	session, ok := r.sessionData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &session, nil
}
func (r *mockUserRepo) DeleteSession(id string) error {
	_, ok := r.sessionData[id]
	if !ok {
		return r.NotFoundErr()
	}
	delete(r.sessionData, id)
	return nil
}
func (r *mockUserRepo) DeleteAllUserSessions(userID string) error {
	for _, s := range r.sessionData {
		if s.UserID == &userID {
			delete(r.sessionData, userID)
		}
	}
	return nil
}

func (r *mockUserRepo) CreateUser(user models.User) (*models.User, error) {
	id := genericUUID
	f := false
	user.ID = &id
	user.Verified = &f
	r.userData[id] = user
	resultUser, ok := r.userData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &resultUser, nil
}
func (r *mockUserRepo) FindUsers() (*[]models.User, error) {
	var u []models.User
	for _, v := range r.userData {
		u = append(u, v)
	}
	return &u, nil
}
func (r *mockUserRepo) FindUserByID(id string) (*models.User, error) {
	user, ok := r.userData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	return &user, nil
}
func (r *mockUserRepo) UpdateUser(id string, user models.User) (*models.User, error) {
	var resultUser models.User
	_, ok := r.userData[id]
	if !ok {
		return nil, r.NotFoundErr()
	}
	resultUser = r.userData[id]
	if user.Username != nil {
		resultUser.Username = user.Username
		r.userData[id] = resultUser
	}
	if user.Password != nil {
		resultUser.Password = user.Password
		r.userData[id] = resultUser
	}
	resultUser = r.userData[id]
	return &resultUser, nil
}
func (r *mockUserRepo) DeleteUser(id string) error {
	_, ok := r.userData[id]
	if !ok {
		return r.NotFoundErr()
	}
	delete(r.userData, id)
	return nil
}

func (r *mockUserRepo) NotFoundErr() error { return r.notFoundErr }
func (r *mockUserRepo) ConflictErr() error { return r.conflictErr }

func (a *mockUserAuthorizer) IdentifyUser(email, username *string) (*models.User, error) {
	for _, v := range a.repository.userData {
		if email != nil {
			if *v.Email == *email {
				return &v, nil
			}
		}
		if username != nil {
			if *v.Username == *username {
				return &v, nil
			}
		}
	}
	return nil, mu.notFoundErr
}

func (a *mockUserAuthorizer) CreatePasswordHash(password *string) (*string, error) {
	return password, nil
}

func (a *mockUserAuthorizer) CheckPassword(password, userID *string) error {
	u, err := a.repository.FindUserByID(*userID)
	if err != nil {
		return err
	}
	if *u.Password != *password {
		return errors.New("Passwords do not match!")
	}
	return nil
}

func TestPostRole(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/roles", strings.NewReader(roleJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.PostRole(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, roleJSONExpected, rec.Body.String())
	}
}

func TestGetRoleByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/roles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetRoleById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, roleJSONExpected, rec.Body.String())
	}
}

func TestGetRoleByIDNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/roles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93ea2872-7da0-49ad-9ff6-a02a99bc3c90")
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetRoleById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestPatchRole(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/roles", strings.NewReader(roleJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.PatchRole(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, roleJSONUpdateExpected, rec.Body.String())
	}
}

func TestPatchRoleNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/roles", strings.NewReader(roleJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93ea2872-7da0-49ad-9ff6-a02a99bc3c90")
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.PatchRole(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestPatchRoleUnprocessable(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/roles", strings.NewReader(roleJSONUpdateInvalid))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.PatchRole(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestPatchRoleConflict(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/roles", strings.NewReader(roleJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.PatchRole(c)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, conflictResponse, rec.Body.String())
	}
}

func TestDeleteRole(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/roles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetRoleById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteRoleUnprocesable(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/roles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("9bc3c90")
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetRoleById(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDeleteRoleNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/roles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93ea2872-7da0-49ad-9ff6-a02a99bc3c90")
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetRoleById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestRegister(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/register", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("password", userPassword)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, userJSONExpected, rec.Body.String())
		s := rec.Result().Cookies()
		assert.Equal(t, genericUUID, s[0].Value)
	}
}

func TestVerify(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/verify", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Sessionid", genericUUID)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.Verify(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, verifiedResponse, rec.Body.String())
	}
}

func TestLogoutHeader(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Sessionid", genericUUID)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, logoutResponse, rec.Body.String())
	}
}

func TestLoginEmail(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/login", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Email", userEmail)
	req.Header.Set("Password", userPassword)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, loginResponse, rec.Body.String())
		s := rec.Result().Cookies()
		assert.Equal(t, genericUUID, s[0].Value)
	}
}

func TestLogoutCookie(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(&http.Cookie{Name: "sessionID", Value: genericUUID})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, logoutResponse, rec.Body.String())
	}
}

func TestLoginUsername(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/login", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Username", userUsername)
	req.Header.Set("password", userPassword)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, loginResponse, rec.Body.String())
		s := rec.Result().Cookies()
		assert.Equal(t, genericUUID, s[0].Value)
	}
}

func TestResetPassword(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/reset-password", strings.NewReader(roleJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("password", userPasswordReset)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.ResetPassword(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, passwordResetResponse, rec.Body.String())
	}
}

func TestLoginIncorrectPassword(t *testing.T) { // Password is changed by this time
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/game-hangar/v1/login", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Username", userUsername)
	req.Header.Set("password", userPassword)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, passwordIncorrectResponse, rec.Body.String())
	}
}

func TestGetUsers(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSONExpectedMany, rec.Body.String())
	}
}

func TestGetUserByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetUserById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSONExpected, rec.Body.String())
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/game-hangar/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93ea2872-7da0-49ad-9ff6-a02a99bc3c90")
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetUserById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFoundResponse, rec.Body.String())
	}
}

func TestPatchUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/game-hangar/v1/users", strings.NewReader(userJSONUpdate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.PatchUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSONUpdateExpected, rec.Body.String())
	}
}

func TestDeleteUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(genericUUID)
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetUserById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteUserUnprocesable(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("9bc3c90")
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetUserById(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/game-hangar/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("93ea2872-7da0-49ad-9ff6-a02a99bc3c90")
	h := &UserHandler{logger: e.Logger, validator: v, repository: &mu, userAuthorizer: &au}

	// Assertions
	if assert.NoError(t, h.GetUserById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
