package routes

import (
	"gamehangar/internal/delivery/http/v1/handlers"

	"github.com/labstack/echo/v4"
)

type DemoRoutes struct {
	handler *handlers.DemoHandler
}

func NewDemoRoutes(h *handlers.DemoHandler) *DemoRoutes {
	return &DemoRoutes{
		handler: h,
	}
}

func (r *DemoRoutes) InitRoutes(e *echo.Echo) {
	demoGroup := e.Group("/game-hangar/v1/demos")

	protectedDemoGroup := demoGroup.Group("")

	protectedDemoGroup.POST("", r.handler.PostDemo)
	demoGroup.GET("/:id", r.handler.GetDemoById)
	demoGroup.GET("", r.handler.GetDemos)
	protectedDemoGroup.PATCH("/:id", r.handler.PatchDemo)
	protectedDemoGroup.DELETE("/:id", r.handler.DeleteDemo)
}
