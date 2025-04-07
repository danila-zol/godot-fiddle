package main

import "github.com/labstack/echo/v4"

func (app *application) routes(e *echo.Echo) *echo.Router {
	router := echo.NewRouter(e)

	v1 := e.Group("/v1")

	userGroup := v1.Group("/user")
	roleGroup := v1.Group("/role")

	demoGroup := v1.Group("/demo")
	assetGroup := v1.Group("/asset")

	topicGroup := v1.Group("/topic")
	threadGroup := v1.Group("/thread")
	messageGroup := v1.Group("/message")

	e.POST("/register", func(c echo.Context) error { return nil })
	e.POST("/login", func(c echo.Context) error { return nil })

	return router
}
