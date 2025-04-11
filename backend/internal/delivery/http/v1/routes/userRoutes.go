package routes

import (
	"gamehangar/internal/delivery/http/v1/handlers"

	"github.com/labstack/echo/v4"
)

type UserRoutes struct {
	handler *handlers.UserHandler
}

func NewUserRoutes(h *handlers.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: h,
	}
}

func (r *UserRoutes) InitRoutes(e *echo.Echo) {
	userGroup := e.Group("/v1/user")

	protectedUserGroup := userGroup.Group("/protected")

	protectedUserGroup.POST("", r.handler.PostUser)
	userGroup.GET("/:id", r.handler.GetUserById)
	userGroup.GET("", r.handler.GetUsers)
	protectedUserGroup.PATCH("/:id", r.handler.PatchUser)
	protectedUserGroup.DELETE("/:id", r.handler.DeleteUser)

	roleGroup := e.Group("/v1/role")

	protectedRoleGroup := roleGroup.Group("/protected")

	protectedRoleGroup.POST("", r.handler.PostRole)
	roleGroup.GET("/:id", r.handler.GetRoleById)
	protectedRoleGroup.PATCH("/:id", r.handler.PatchRole)
	protectedRoleGroup.DELETE("/:id", r.handler.DeleteRole)

	sessionGroup := e.Group("/v1/session")

	sessionGroup.POST("", r.handler.PostSession)
	sessionGroup.GET("/:id", r.handler.GetSessionById)
	sessionGroup.DELETE("/:id", r.handler.DeleteSession)
}
