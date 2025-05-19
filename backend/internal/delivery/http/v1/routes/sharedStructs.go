package routes

import "github.com/labstack/echo/v4"

type Authorizer interface {
	CheckPermissions(c echo.Context, user string) (bool, error)
}
