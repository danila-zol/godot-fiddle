package handlers

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"time"

	_ "gamehangar/docs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserIdentifier interface {
	IdentifyUser(email, username *string) (user *models.User, err error)
}

type UserHandler struct {
	logger         echo.Logger
	repository     UserRepository
	userIdentifier UserIdentifier
}

func NewUserHandler(e *echo.Echo, repo UserRepository, r UserIdentifier) *UserHandler {
	return &UserHandler{
		logger:         e.Logger,
		repository:     repo,
		userIdentifier: r,
	}
}

// @Summary	Fetches a user by its ID.
// @Tags		Users
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get User of ID"
// @Success	200	{object}	ResponseHTTP{data=models.User}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/users/{id} [get]
func (h *UserHandler) GetUserById(c echo.Context) error {
	id := c.Param("id")

	user, err := h.repository.FindUserByID(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Users not found! %s", err)
			return c.String(http.StatusNotFound, "Error: User not found!")
		}
		h.logger.Printf("Error in FindUserByID repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateUser repository")
	}

	return c.JSON(http.StatusOK, &user)
}

// @Summary	Fetches all users.
// @Tags		Users
// @Produce	application/json
// @Success	200	{object}	ResponseHTTP{data=[]models.User}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/users [get]
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.repository.FindUsers()
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: User not found! %s", err)
			return c.String(http.StatusNotFound, "Error: User not found!")
		}
		h.logger.Printf("Error in FindUsers operation: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindUsers operation")
	}

	return c.JSON(http.StatusOK, &users)
}

// @Summary	Updates a user.
// @Tags		Users
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		string		true	"Update User of ID"
// @Param		User	body		models.User	true	"Update User"
// @Success	200		{object}	ResponseHTTP{data=models.User}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/users/{id} [patch]
func (h *UserHandler) PatchUser(c echo.Context) error {
	var user models.User

	id := c.Param("id")

	err := c.Bind(&user)
	if err != nil {
		h.logger.Printf("Error in PatchUser handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchUser handler")
	}

	updUser, err := h.repository.UpdateUser(id, user)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: User not found! %s", err)
			return c.String(http.StatusNotFound, "Error: User not found!")
		}
		h.logger.Printf("Error in UpdateUser repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateUser repository")
	}

	return c.JSON(http.StatusOK, &updUser)
}

// @Summary	Deletes the specified user.
// @Tags		Users
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete User of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	err := h.repository.DeleteUser(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: User not found! %s", err)
			return c.String(http.StatusNotFound, "Error: User not found!")
		}
		h.logger.Printf("Error in DeleteUser repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteUser repository")
	}

	return c.String(http.StatusOK, "User successfully deleted!")
}

// @Summary	Creates a new role.
// @Tags		Roles
// @Accept		application/json
// @Produce	application/json
// @Param		Role	body		models.Role	true	"Create Role"
// @Success	200		{object}	ResponseHTTP{data=models.Role}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/roles [post]
func (h *UserHandler) PostRole(c echo.Context) error {
	var role models.Role

	err := c.Bind(&role)
	if err != nil {
		h.logger.Printf("Error in PostRole handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostRole handler")
	}

	if role.ID == nil {
		roleID := uuid.NewString()
		role.ID = &roleID
	}

	newRole, err := h.repository.CreateRole(role)
	if err != nil {
		h.logger.Printf("Error in CreateRole repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateRole repository")
	}

	return c.JSON(http.StatusOK, &newRole)
}

// @Summary	Fetches a role by its ID.
// @Tags		Roles
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Role of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Role}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/roles/{id} [get]
func (h *UserHandler) GetRoleById(c echo.Context) error {
	id := c.Param("id")

	role, err := h.repository.FindRoleByID(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Role not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Role not found!")
		}
		h.logger.Printf("Error in FindRoleByID repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateRole repository")
	}

	return c.JSON(http.StatusOK, &role)
}

// @Summary	Updates a role.
// @Tags		Roles
// @Accept		application/json
// @Produce	application/json
// @Param		id		path		string		true	"Update Role of ID"
// @Param		Role	body		models.Role	true	"Update Role"
// @Success	200		{object}	ResponseHTTP{data=models.Role}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/roles/{id} [patch]
func (h *UserHandler) PatchRole(c echo.Context) error {
	var role models.Role

	id := c.Param("id")

	err := c.Bind(&role)
	if err != nil {
		h.logger.Printf("Error in PatchRole handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PatchRole handler")
	}

	updRole, err := h.repository.UpdateRole(id, role)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Role not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Role not found!")
		}
		h.logger.Printf("Error in UpdateRole repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateRole repository")
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
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/roles/{id} [delete]
func (h *UserHandler) DeleteRole(c echo.Context) error {
	id := c.Param("id")

	err := h.repository.DeleteRole(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Role not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Role not found!")
		}
		h.logger.Printf("Error in DeleteRole repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteRole repository")
	}

	return c.String(http.StatusOK, "Role successfully deleted!")
}

// @Summary	Registers a new user and creates a session.
// @Tags		Login
// @Accept		application/json
// @Produce	application/json
// @Param		User	body		models.User	true	"Create User"
// @Success	200		{object}	ResponseHTTP{data=models.User}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		h.logger.Printf("Error in PostUser handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostUser handler")
	}

	if user.CreatedAt == nil {
		currentTime := time.Now()
		user.CreatedAt = &currentTime
	}
	if user.Karma == nil {
		zero := 0
		user.Karma = &zero
	}
	if user.Verified == nil {
		f := false
		user.Verified = &f
	} // TODO: Fix nil values for Postgres
	// TODO: Create Salt for user's password

	newUser, err := h.repository.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error in CreateUser repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateUser repository")
	}

	session, err := h.repository.CreateSession(models.Session{UserID: newUser.ID})
	if err != nil {
		h.logger.Printf("Error creating session: %s", err)
		return c.String(http.StatusInternalServerError, "Error creating session")
	}
	c.SetCookie(&http.Cookie{
		Name:     "sessionID",
		Value:    *session.ID,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		SameSite: 3,
	})
	h.logger.Printf("Cookie value: ", *session.ID) // DEBUG

	return c.JSON(http.StatusOK, &newUser)
}

// @Summary	Verifies the authenticated User.
// @Tags		Login
// @Accept		text/plain
// @Produce	text/plain
// @Security ApiSessionCookie
// @param sessionID header string true "Session ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/verify [get]
func (h *UserHandler) Verify(c echo.Context) error {
	cookie, err := c.Cookie("sessionID")
	if err != nil {
		h.logger.Printf("Error reading cookie: %s", err)
		return c.String(http.StatusInternalServerError, "Error reading cookie")
	}

	session, err := h.repository.FindSessionByID(cookie.Value)
	if err != nil {
		h.logger.Printf("Error in FindSessionByID repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in FindSessionByID repository")
	}

	t := true
	_, err = h.repository.UpdateUser(*session.UserID, models.User{Verified: &t})
	if err != nil {
		h.logger.Printf("Error in UpdateUser repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateUser repository")
	}

	return c.String(http.StatusOK, "User verified")
}

// @Summary	Logs the User in and creates a new Session.
// @Tags		Login
// @Accept		application/json
// @Produce	text/plain
// @Param		LoginForm	body		models.LoginForm	true	"Log in"
// @Success	200		{object}	ResponseHTTP{}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var loginForm models.LoginForm

	err := c.Bind(&loginForm)
	if err != nil {
		h.logger.Printf("Error in PostSession handler: %s", err)
		return c.String(http.StatusBadRequest, "Error in PostSession handler")
	}

	user, err := h.userIdentifier.IdentifyUser(loginForm.Email, loginForm.Username)
	if err != nil {
		h.logger.Printf("Error identifying user: %s", err)
		return c.String(http.StatusInternalServerError, "Error identifying user")
	}

	// TODO: Service to check password

	session, err := h.repository.CreateSession(models.Session{UserID: user.ID})
	if err != nil {
		h.logger.Printf("Error in CreateSession repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateSession repository")
	}

	c.SetCookie(&http.Cookie{
		Name:     "sessionID",
		Value:    *session.ID,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		SameSite: 3,
	})
	h.logger.Printf("Cookie value: ", *session.ID) // DEBUG

	return c.String(http.StatusOK, "Login successful")
}

// @Summary	Refreshes Session token.
// @Tags		Login
// @Accept		text/plain
// @Produce	text/plain
// @Security ApiSessionCookie
// @param sessionID header string true "Session ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/refresh [get]
func (h *UserHandler) RefreshSession(c echo.Context) error {
	cookie, err := c.Cookie("sessionID")
	if err != nil {
		h.logger.Printf("Error reading cookie: %s", err)
		return c.String(http.StatusBadRequest, "Error reading cookie")
	}

	/* No need to create a new session, right?
	id, err := c.Cookie("sessionID")
	if err != nil {
		h.logger.Printf("Error finding session in cookie: %s", err)
		return c.String(http.StatusInternalServerError, "Error finding session in cookie")
	}
	*/

	c.SetCookie(&http.Cookie{
		Name:     "sessionID",
		Value:    cookie.Value,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		SameSite: 3,
	})

	return c.String(http.StatusOK, "Session refreshed")
}

// @Summary	Invalidates and deletes user's Session.
// @Tags		Login
// @Accept		text/plain
// @Produce	text/plain
// @Security ApiSessionCookie
// @param sessionID header string true "Session ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/logout [delete]
func (h *UserHandler) Logout(c echo.Context) error {
	sessionID, err := c.Cookie("sessionID")
	if err != nil {
		h.logger.Printf("Error reading cookie: %s", err)
		return c.String(http.StatusBadRequest, "Error reading cookie")
	}

	err = h.repository.DeleteSession(sessionID.Value)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Sessions not found! %s", err)
			return c.String(http.StatusNotFound, "Error: Session not found!")
		}
		h.logger.Printf("Error in DeleteSession repository: %s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteSession repository")
	}

	c.SetCookie(&http.Cookie{
		Name:     "sessionID",
		Value:    "",
		Expires:  time.Now().Add(-1),
		HttpOnly: true,
		SameSite: 3,
	})

	return c.String(http.StatusOK, "Session successfully deleted!")
}
