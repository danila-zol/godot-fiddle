package handlers

import (
	"gamehangar/internal/domain/models"
	"net/http"
	"time"

	_ "gamehangar/docs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	logger     echo.Logger
	repository UserRepository
}

func NewUserHandler(e *echo.Echo, repo UserRepository) *UserHandler {
	return &UserHandler{
		logger:     e.Logger,
		repository: repo,
	}
}

// @Summary	Creates a new user.
// @Tags		Users
// @Accept		application/json
// @Produce	application/json
// @Param		User	body		models.User	true	"Create User"
// @Success	200		{object}	ResponseHTTP{data=models.User}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/user/ [post]
func (h *UserHandler) PostUser(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		h.logger.Printf("Error in PostUser handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PostUser handler")
	}

	if user.ID == nil {
		userID := uuid.NewString()
		user.ID = &userID
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
	}
	// TODO: Create Salt for user's password

	newUser, err := h.repository.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error in CreateUser repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateUser repository")
	}

	return c.JSON(http.StatusOK, &newUser)
}

// @Summary	Fetches a user by its ID.
// @Tags		Users
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get User of ID"
// @Success	200	{object}	ResponseHTTP{data=models.User}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/user/{id} [get]
func (h *UserHandler) GetUserById(c echo.Context) error {
	id := c.Param("id")

	user, err := h.repository.FindUserByID(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Users not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: User not found!")
		}
		h.logger.Printf("Error in FindUserByID repository \n%s", err)
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
// @Router		/v1/user/ [get]
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.repository.FindUsers()
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Users not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in FindUsers operation \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in FindUsers operation")
	}

	return c.JSON(http.StatusOK, &users)
}

// @Summary	Updates a user.
// @Tags		Users
// @Accept		application/json
// @Produce	application/json
// @Param		id	path		string	true	"Update User of ID"
// @Param		User	body		models.User	true	"Update User"
// @Success	200		{object}	ResponseHTTP{data=models.User}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/user/protected/{id} [patch]
func (h *UserHandler) PatchUser(c echo.Context) error {
	var user models.User

	id := c.Param("id")

	err := c.Bind(&user)
	if err != nil {
		h.logger.Printf("Error in PatchUser handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PatchUser handler")
	}

	updUser, err := h.repository.UpdateUser(id, user)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Users not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in UpdateUser repository \n%s", err)
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
// @Router		/v1/user/protected/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	err := h.repository.DeleteUser(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Users not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in DeleteUser repository \n%s", err)
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
// @Router		/v1/role/ [post]
func (h *UserHandler) PostRole(c echo.Context) error {
	var role models.Role

	err := c.Bind(&role)
	if err != nil {
		h.logger.Printf("Error in PostRole handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PostRole handler")
	}

	if role.ID == nil {
		roleID := uuid.NewString()
		role.ID = &roleID
	}

	newRole, err := h.repository.CreateRole(role)
	if err != nil {
		h.logger.Printf("Error in CreateRole repository \n%s", err)
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
// @Router		/v1/role/{id} [get]
func (h *UserHandler) GetRoleById(c echo.Context) error {
	id := c.Param("id")

	role, err := h.repository.FindRoleByID(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Roles not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: Role not found!")
		}
		h.logger.Printf("Error in FindRoleByID repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateRole repository")
	}

	return c.JSON(http.StatusOK, &role)
}

// @Summary	Updates a role.
// @Tags		Roles
// @Accept		application/json
// @Produce	application/json
// @Param		id	path		string	true	"Update Role of ID"
// @Param		Role	body		models.Role	true	"Update Role"
// @Success	200		{object}	ResponseHTTP{data=models.Role}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/role/protected/{id} [patch]
func (h *UserHandler) PatchRole(c echo.Context) error {
	var role models.Role

	id := c.Param("id")

	err := c.Bind(&role)
	if err != nil {
		h.logger.Printf("Error in PatchRole handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PatchRole handler")
	}

	updRole, err := h.repository.UpdateRole(id, role)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Roles not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in UpdateRole repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateRole repository")
	}

	return c.JSON(http.StatusOK, &updRole)
}

// @Summary	Deletes the specified role.
// @Tags		Roles
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Role of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/role/protected/{id} [delete]
func (h *UserHandler) DeleteRole(c echo.Context) error {
	id := c.Param("id")

	err := h.repository.DeleteRole(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Roles not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in DeleteRole repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteRole repository")
	}

	return c.String(http.StatusOK, "Role successfully deleted!")
}

// @Summary	Creates a new session.
// @Tags		Sessions
// @Accept		application/json
// @Produce	application/json
// @Param		Session	body		models.Session	true	"Create Session"
// @Success	200		{object}	ResponseHTTP{data=models.Session}
// @Failure	400		{object}	ResponseHTTP{}
// @Failure	500		{object}	ResponseHTTP{}
// @Router		/v1/session/ [post]
func (h *UserHandler) PostSession(c echo.Context) error {
	var session models.Session

	err := c.Bind(&session)
	if err != nil {
		h.logger.Printf("Error in PostSession handler \n%s", err)
		return c.String(http.StatusBadRequest, "Error in PostSession handler")
	}

	if session.ID == nil {
		sessionID := uuid.NewString()
		session.ID = &sessionID
	}

	newSession, err := h.repository.CreateSession(session)
	if err != nil {
		h.logger.Printf("Error in CreateSession repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in CreateSession repository")
	}

	return c.JSON(http.StatusOK, &newSession)
}

// @Summary	Fetches a session by its ID.
// @Tags		Sessions
// @Accept		text/plain
// @Produce	application/json
// @Param		id	path		string	true	"Get Session of ID"
// @Success	200	{object}	ResponseHTTP{data=models.Session}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/session/{id} [get]
func (h *UserHandler) GetSessionById(c echo.Context) error {
	id := c.Param("id")

	session, err := h.repository.FindSessionByID(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Sessions not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: Session not found!")
		}
		h.logger.Printf("Error in FindSessionByID repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in UpdateSession repository")
	}

	return c.JSON(http.StatusOK, &session)
}

// @Summary	Deletes the specified session.
// @Tags		Sessions
// @Accept		text/plain
// @Produce	text/plain
// @Param		id	path		string	true	"Delete Session of ID"
// @Success	200	{object}	ResponseHTTP{}
// @Failure	400	{object}	ResponseHTTP{}
// @Failure	500	{object}	ResponseHTTP{}
// @Router		/v1/session/{id} [delete]
func (h *UserHandler) DeleteSession(c echo.Context) error {
	id := c.Param("id")

	err := h.repository.DeleteSession(id)
	if err != nil {
		if err == h.repository.NotFoundErr() {
			h.logger.Printf("Error: Sessions not found!\n%s", err)
			return c.String(http.StatusNotFound, "Error: Asset not found!")
		}
		h.logger.Printf("Error in DeleteSession repository \n%s", err)
		return c.String(http.StatusInternalServerError, "Error in DeleteSession repository")
	}

	return c.String(http.StatusOK, "Session successfully deleted!")
}
