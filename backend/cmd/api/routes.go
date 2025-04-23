package main

import (
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	_ "gamehangar/docs"
)

// Configured to init v1 routes and handlers
func (app *application) routes(e *echo.Echo) *echo.Router {
	router := echo.NewRouter(e)
	skipCSRF := []string{"/game-hangar/v1/register", "/game-hangar/v1/login"}

	e.GET("/game-hangar/docs/*", echoSwagger.WrapHandler)
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper: func(c echo.Context) bool {
			if slices.Contains(skipCSRF, c.Request().URL.Path) {
				return true
			}
			return false
		},
		TokenLookup: "cookie:_csrf",
	}))

	return router
}
