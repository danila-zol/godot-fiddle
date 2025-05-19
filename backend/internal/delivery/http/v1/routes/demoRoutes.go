package routes

import (
	"gamehangar/internal/delivery/http/v1/handlers"

	casbin_mw "github.com/labstack/echo-contrib/casbin"
	"github.com/labstack/echo/v4"
)

type DemoRoutes struct {
	handler    *handlers.DemoHandler
	authorizer Authorizer
}

func NewDemoRoutes(h *handlers.DemoHandler, a Authorizer) *DemoRoutes {
	return &DemoRoutes{
		handler:    h,
		authorizer: a,
	}
}

func (r *DemoRoutes) InitRoutes(e *echo.Echo) {
	demoGroup := e.Group("/game-hangar/v1/demos")

	protectedDemoGroup := demoGroup.Group("")
	protectedDemoGroup.Use(casbin_mw.MiddlewareWithConfig(casbin_mw.Config{
		EnforceHandler: r.authorizer.CheckPermissions,
	}))

	protectedDemoGroup.POST("", r.handler.PostDemo)
	demoGroup.GET("/:id", r.handler.GetDemoById)
	demoGroup.GET("", r.handler.GetDemos)
	protectedDemoGroup.PATCH("/:id", r.handler.PatchDemo)
	protectedDemoGroup.DELETE("/:id", r.handler.DeleteDemo)
}
