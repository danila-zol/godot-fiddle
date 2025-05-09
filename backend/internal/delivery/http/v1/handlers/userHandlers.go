package handlers

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"strconv"
	"time"

	_ "gamehangar/docs"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserAuthorizer interface {
	IdentifyUser(email, username *string) (user *models.User, err error)
	CreatePasswordHash(password *string) (hash *string, err error)
	CheckPassword(password *string, userID uuid.UUID) (err error)
}

type UserHandler struct {
	logger         echo.Logger
	repository     UserRepository
	validator      *validator.Validate
	userAuthorizer UserAuthorizer
}

func NewUserHandler(e *echo.Echo, repo UserRepository, v *validator.Validate, a UserAuthorizer) *UserHandler {
	return &UserHandler{
		logger:         e.Logger,
		repository:     repo,
		validator:      v,
		userAuthorizer: a,
	}
}

// @Summary	Fetches a user by its ID.
// @Tags		Users
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get User of ID"
// @Success	200	{object}	models.User
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/users/{id} [get]
func (h *UserHandler) GetUserById(c echo.Context) error {
	id := c.Param("id")
	err := h.validator.Var(id, "required,uuid4")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in GetUserByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	userID, _ := uuid.Parse(id)
	user, err := h.repository.FindUserByID(userID)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindUserByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &user)
}

// @Summary	Fetches all users.
// @Tags		Users
// @Produce	application/json
// @Param		q	query		[]string	false	"Keyword Query"
// @Param		l	query		int	false	"Record number limit"
// @Success	200	{object}	models.User
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/users [get]
func (h *UserHandler) GetUsers(c echo.Context) error {
	var err error
	var limit uint64
	var users *[]models.User

	l := c.Request().URL.Query()["l"]
	if l != nil {
		err = h.validator.Var(l[0], "omitnil,number,min=0")
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in GetUsers repository: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnprocessableEntity, &e)
		}
		limit, err = strconv.ParseUint(l[0], 10, 64)
	}
	tags := c.Request().URL.Query()["q"]

	users, err = h.repository.FindUsers(tags, limit)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindUsers repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &users)
}

// @Summary	Updates a user.
// @Tags		Users
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		string		true	"Update User of ID"
// @Param		User	body		models.User	true	"Update User"
// @Success	200		{object}	models.User
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/users/{id} [patch]
func (h *UserHandler) PatchUser(c echo.Context) error {
	var user models.User

	id := c.Param("id")
	err := h.validator.Var(id, "required,uuid4")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchUser handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	err = c.Bind(&user)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PatchUser handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}

	err = h.validator.Struct(&user)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchUser handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	userID, _ := uuid.Parse(id)
	updUser, err := h.repository.UpdateUser(userID, user)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in UpdateUser repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &updUser)
}

// @Summary	Deletes the specified user.
// @Tags		Users
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete User of ID"
// @Success	200	{string}	string
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.validator.Var(id, "required,uuid4")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in DeleteUser handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	userID, _ := uuid.Parse(id)
	err = h.repository.DeleteUser(userID)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in DeleteUser repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "User successfully deleted!")
}

// @Summary	Creates a new role.
// @Tags		Roles
// @Accept		application/json
// @Produce	application/json
// @Param		Role	body		models.Role	true	"Create Role"
// @Success	201	{object}	models.Role
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/roles [post]
func (h *UserHandler) PostRole(c echo.Context) error {
	var role models.Role

	err := c.Bind(&role)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PostRole handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	role.Method = "POST"

	err = h.validator.Struct(&role)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PostRole handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	newRole, err := h.repository.CreateRole(role)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateRole repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusCreated, &newRole)
}

// @Summary	Fetches a role by its ID.
// @Tags		Roles
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Role of ID"
// @Success	200	{object}	models.Role
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/roles/{id} [get]
func (h *UserHandler) GetRoleById(c echo.Context) error {
	id := c.Param("id")
	err := h.validator.Var(id, "required,uuid4")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in GetRoleByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	roleID, _ := uuid.Parse(id)
	role, err := h.repository.FindRoleByID(roleID)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindRoleByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &role)
}

// @Summary	Updates a role.
// @Tags		Roles
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		string		true	"Update Role of ID"
// @Param		Role	body		models.Role	true	"Update Role"
// @Success	200		{object}	models.Role
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	409	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/roles/{id} [patch]
func (h *UserHandler) PatchRole(c echo.Context) error {
	var role models.Role

	id := c.Param("id")
	err := h.validator.Var(id, "required,uuid4")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchRole handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	err = c.Bind(&role)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in PatchRole handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	role.Method = "PATCH"

	err = h.validator.Struct(&role)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in PatchRole handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	roleID, _ := uuid.Parse(id)
	updRole, err := h.repository.UpdateRole(roleID, role)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		} else if err == h.repository.ConflictErr() {
			e := HTTPError{
				Code:    http.StatusConflict,
				Message: "Error: unable to update the record due to an edit conflict, please try again!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusConflict, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in UpdateRole repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.JSON(http.StatusOK, &updRole)
}

// @Summary	Deletes the specified role.
// @Tags		Roles
// @Accept		text/plain
// @Produce	text/plain
// @Security ApiSessionCookie
// @param sessionID header string true "Session ID"
// @Param		id	path		string	true	"Delete Role of ID"
// @Success	200	{string}	string
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/roles/{id} [delete]
func (h *UserHandler) DeleteRole(c echo.Context) error {
	id := c.Param("id")
	err := h.validator.Var(id, "required,uuid4")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in DeleteRole handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	roleID, _ := uuid.Parse(id)
	err = h.repository.DeleteRole(roleID)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in DeleteRole repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "Role successfully deleted!")
}

// @Summary	Registers a new user and creates a session.
// @Tags		Login
// @Accept		application/json
// @Produce	application/json
// @Param		User	body		models.User	true	"Create User"
// @param password header string true "Password"
// @Success	201	{object}	models.User
// @Failure	400	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Error in Register handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	user.Method = "POST"

	err = h.validator.Struct(&user)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in Register handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	passwordSlice, ok := c.Request().Header["Password"]
	if !ok {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "No password provided!",
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	password := &passwordSlice[0]
	err = h.validator.Var(&password, "required,min=8")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in Register handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}
	user.Password, err = h.userAuthorizer.CreatePasswordHash(password)

	newUser, err := h.repository.CreateUser(user)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateUser repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	session, err := h.repository.CreateSession(models.Session{UserID: newUser.ID})
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateSession repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}
	c.SetCookie(&http.Cookie{
		Name:     "sessionID",
		Value:    session.ID.String(),
		Expires:  time.Now().Add(96 * time.Hour),
		HttpOnly: true,
		SameSite: 3,
	})
	h.logger.Printf("Cookie value: ", *session.ID) // DEBUG

	return c.JSON(http.StatusCreated, &newUser)
}

// @Summary	Verifies the authenticated User.
// @Tags		Login
// @Accept		text/plain
// @Produce	text/plain
// @Security ApiSessionCookie
// @param sessionID header string false "Session ID"
// @Success	200	{string}	string
// @Failure	401	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/verify [get]
func (h *UserHandler) Verify(c echo.Context) error {
	var s string

	cookie, err := c.Cookie("sessionID")
	if err != nil {
		sessionSlice, ok := c.Request().Header["Sessionid"]
		if !ok {
			e := HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "No session provided!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnauthorized, &e)
		}
		s = sessionSlice[0]
	} else {
		s = cookie.Value
	}

	sessionID, _ := uuid.Parse(s)
	session, err := h.repository.FindSessionByID(sessionID)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusUnauthorized, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnauthorized, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindSessionByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	t := true
	_, err = h.repository.UpdateUser(*session.UserID, models.User{Verified: &t})
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in UpdateUser repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "User verified")
}

// @Summary	Resets User password and deletes all their Sessions
// @Tags		Login
// @Accept		text/plain
// @Produce	text/plain
// @Security ApiSessionCookie
// @param password header string true "New Password"
// @param id path string true "User ID"
// @Success	200	{string}	string
// @Failure	400	{object}	HTTPError
// @Failure	422	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/reset-password/{id} [patch]
func (h *UserHandler) ResetPassword(c echo.Context) error {
	passwordSlice, ok := c.Request().Header["Password"]
	if !ok {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "No password provided!",
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	password := &passwordSlice[0]

	err := h.validator.Var(&password, "required,min=8")
	if err != nil {
		e := HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: "Error in ResetPassword handler: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnprocessableEntity, &e)
	}

	password, err = h.userAuthorizer.CreatePasswordHash(password)

	sessionID, _ := uuid.Parse(c.Param("id"))
	user, err := h.repository.UpdateUser(sessionID, models.User{Password: password})
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in UpdateUser repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	err = h.repository.DeleteAllUserSessions(*user.ID)
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in DeleteAllUserSessions repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	return c.String(http.StatusOK, "User password reset!")
}

// @Summary	Logs the User in and creates a new Session.
// @Tags		Login
// @Accept		application/json
// @Produce	text/plain
// @param email header string false "Email"
// @param username header string false "Username"
// @param password header string true "Password"
// @Success	200		{string}	string
// @Failure	400		{object}	HTTPError
// @Failure	500		{object}	HTTPError
// @Router		/v1/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var username, email, password string

	passwordSlice, ok := c.Request().Header["Password"]
	if !ok {
		e := HTTPError{
			Code:    http.StatusBadRequest,
			Message: "No password provided!",
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusBadRequest, &e)
	}
	password = passwordSlice[0]

	if _, ok = c.Request().Header["Username"]; ok {
		username = c.Request().Header["Username"][0]
		err := h.validator.Var(username, "omitempty,max=90")
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in Login handler: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnprocessableEntity, &e)
		}
	} else if _, ok = c.Request().Header["Email"]; ok {
		email = c.Request().Header["Email"][0]
		err := h.validator.Var(email, "omitempty,email,max=50")
		if err != nil {
			e := HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "Error in Login handler: " + err.Error(),
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnprocessableEntity, &e)
		}
	} else {
		e := HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "No identifiers provided!",
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnauthorized, &e)
	}

	user, err := h.userAuthorizer.IdentifyUser(&email, &username)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusNotFound, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in IdentifyUser service: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	if h.userAuthorizer.CheckPassword(&password, *user.ID) != nil {
		e := HTTPError{Code: http.StatusUnauthorized, Message: "Password incorrect!"}
		h.logger.Print(&e)
		return c.JSON(http.StatusUnauthorized, &e)
	}

	session, err := h.repository.CreateSession(models.Session{UserID: user.ID})
	if err != nil {
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in CreateSession repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	c.SetCookie(&http.Cookie{
		Name:     "sessionID",
		Value:    session.ID.String(),
		Expires:  time.Now().Add(96 * time.Hour),
		HttpOnly: true,
		SameSite: 3,
	})
	h.logger.Printf("Cookie value: ", *session.ID) // DEBUG

	return c.String(http.StatusOK, "Login successful")
}

// @Summary	Invalidates and deletes the specified session.
// @Tags		Login
// @Accept		text/plain
// @Produce	text/plain
// @Security ApiSessionCookie
// @param sessionID header string false "Session ID"
// @Param		id		path		string		true	"Session to invalidate"
// @Success	200	{string}	string
// @Failure	400	{object}	HTTPError
// @Failure	401	{object}	HTTPError
// @Failure	403	{object}	HTTPError
// @Failure	404	{object}	HTTPError
// @Failure	500	{object}	HTTPError
// @Router		/v1/logout/{id} [delete]
func (h *UserHandler) Logout(c echo.Context) error {
	var s string

	cookie, err := c.Cookie("sessionID")
	if err != nil {
		sessionSlice, ok := c.Request().Header["Sessionid"]
		if !ok {
			e := HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "No session provided!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnauthorized, &e)
		}
		s = sessionSlice[0]
	} else {
		s = cookie.Value
	}

	sessionID, _ := uuid.Parse(s)
	currentSession, err := h.repository.FindSessionByID(sessionID)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusUnauthorized, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnauthorized, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindSessionByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	reqSessionID, _ := uuid.Parse(c.Param("id"))
	requestedSession, err := h.repository.FindSessionByID(reqSessionID)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{Code: http.StatusUnauthorized, Message: "Not Found!"}
			h.logger.Print(&e)
			return c.JSON(http.StatusUnauthorized, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in FindSessionByID repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}
	if *requestedSession.UserID != *currentSession.UserID { // TODO: Skip for admin privileges
		e := HTTPError{
			Code:    http.StatusForbidden,
			Message: "Cannot log out other user!",
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusForbidden, &e)
	}

	err = h.repository.DeleteSession(*requestedSession.ID)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			e := HTTPError{
				Code:    http.StatusNotFound,
				Message: "Not Found!",
			}
			h.logger.Print(&e)
			return c.JSON(http.StatusNotFound, &e)
		}
		e := HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error in DeleteSession repository: " + err.Error(),
		}
		h.logger.Print(&e)
		return c.JSON(http.StatusInternalServerError, &e)
	}

	// Invalidate cookie if the client invalidates the session they are logged in as
	if *requestedSession.ID == *currentSession.ID {
		c.SetCookie(&http.Cookie{
			Name:     "sessionID",
			Value:    "",
			Expires:  time.Now().Add(-1),
			HttpOnly: true,
			SameSite: 3,
		})
	}

	return c.String(http.StatusOK, "Session successfully deleted!")
}
