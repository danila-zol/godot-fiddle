package main

import (
	"gamehangar/internal/config"
	"gamehangar/internal/database/psql"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type appConfig struct {
	port int
	env  string
}

type application struct {
	appConfig appConfig
	logger    echo.Logger
	validator echo.Validator
}

func (a *application) getEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		a.logger.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	var cfg appConfig
	e := echo.New()

	app := &application{
		appConfig: cfg,
		logger:    e.Logger,
		validator: e.Validator,
	}
	app.getEnv()

	databaseConfig, err := config.NewConfig(psql.MigrationFiles, os.Getenv("PSQL_MIGRATE_FILE_DIR"))
	if err != nil {
		app.logger.Fatalf("Error loading PSQL database Config: %v", err)
	}
	databaseClient, err := psql.NewDatabaseClient(os.Getenv("PSQL_CONNSTRING"), databaseConfig).Setup()
	if err != nil {
		app.logger.Fatalf("Error creating new DatabaseClient: %v", err)
	}
	println(databaseClient)
}
