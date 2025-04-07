package main

import "github.com/labstack/echo/v4"

type config struct {
	port int
	env  string
}

type application struct {
	config    config
	logger    any
	validator any
}

func main() {
	var cfg config
	e := echo.New()

	_ = &application{
		config:    cfg,
		logger:    e.Logger,
		validator: e.Validator,
	}
}
