package routes

import (
	"gamehangar/internal/delivery/http/v1/handlers"

	casbin_mw "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
)

type UserRoutes struct {
	handler    *handlers.UserHandler
	authorizer Authorizer
}

func NewUserRoutes(h *handlers.UserHandler, a Authorizer) *UserRoutes {
	return &UserRoutes{
		handler:    h,
		authorizer: a,
	}
}

func (r *UserRoutes) InitRoutes(e *echo.Echo) {
	userGroup := e.Group("/game-hangar/v1/users")

	protectedUserGroup := userGroup.Group("")
	protectedUserGroup.Use(casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
		EnforceHandler: r.authorizer.CheckPermissions,
	}))

	userGroup.GET("/:id", r.handler.GetUserById)
	userGroup.GET("", r.handler.GetUsers)
	protectedUserGroup.PATCH("/:id", r.handler.PatchUser)
	protectedUserGroup.DELETE("/:id", r.handler.DeleteUser)

	roleGroup := e.Group("/game-hangar/v1/roles")

	protectedRoleGroup := roleGroup.Group("")
	protectedRoleGroup.Use(casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
		EnforceHandler: r.authorizer.CheckPermissions,
	}))

	protectedRoleGroup.POST("", r.handler.PostRole)
	protectedRoleGroup.DELETE("", r.handler.DeleteRole)

	sessionGroup := e.Group("/game-hangar/v1")
	protectedSessionGroup := sessionGroup.Group("")
	protectedSessionGroup.Use(casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
		EnforceHandler: r.authorizer.CheckPermissions,
	}))

	sessionGroup.POST("/register", r.handler.Register)
	sessionGroup.POST("/login", r.handler.Login)
	sessionGroup.PATCH("/reset-password/:id", r.handler.ResetPassword)
	sessionGroup.GET("/verify", r.handler.Verify)
	protectedSessionGroup.DELETE("/logout/:id", r.handler.Logout)
}
