package main

import (
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "gamehangar/docs"
)

// Configured to init v1 routes and handlers
func (app *application) routes(e *echo.Echo) *echo.Router {
	router := echo.NewRouter(e)

	e.GET("/game-hangar/docs/*", echoSwagger.WrapHandler)

	return router
}
