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
	userGroup := e.Group("/game-hangar/v1/users")

	protectedUserGroup := userGroup.Group("")

	userGroup.GET("/:id", r.handler.GetUserById)
	userGroup.GET("", r.handler.GetUsers)
	protectedUserGroup.PATCH("/:id", r.handler.PatchUser)
	protectedUserGroup.DELETE("/:id", r.handler.DeleteUser)

	roleGroup := e.Group("/game-hangar/v1/roles")

	protectedRoleGroup := roleGroup.Group("")

	protectedRoleGroup.POST("", r.handler.PostRole)
	roleGroup.GET("/:id", r.handler.GetRoleById)
	protectedRoleGroup.PATCH("/:id", r.handler.PatchRole)
	protectedRoleGroup.DELETE("/:id", r.handler.DeleteRole)

	sessionGroup := e.Group("/game-hangar/v1")

	sessionGroup.POST("/register", r.handler.Register)
	sessionGroup.POST("/login", r.handler.Login)
	sessionGroup.GET("/verify", r.handler.Verify)
	sessionGroup.GET("/refresh", r.handler.RefreshSession)
	sessionGroup.DELETE("/logout", r.handler.Logout)
}
