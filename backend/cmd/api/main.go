package main

import (
	"gamehangar/internal/config"
	"gamehangar/internal/database"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// Exclusive to package main since it is a central app config
type appConfig struct {
	host string
	port int
}

type application struct {
	appConfig *appConfig
	logger    echo.Logger
	validator echo.Validator
}

func getEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	e := echo.New()
	getEnv()

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 32)
	if err != nil {
		panic("Could not parse app config: Invalid port")
	}

	cfg := &appConfig{
		host: os.Getenv("HOST"),
		port: int(port),
	}

	app := &application{
		appConfig: cfg,
		logger:    e.Logger,
		validator: e.Validator,
	}

	// Might not look pretty to avoid passing around a pointer, but Psql{} structs are empty
	databaseConfig, err := config.Psql{}.NewConfig(database.MigrationFiles, os.Getenv("PSQL_MIGRATE_ROOT_DIR"))
	if err != nil {
		app.logger.Fatalf("Error loading PSQL database Config: %v", err)
	}
	databaseClient := database.Psql{}.NewDatabaseClient(os.Getenv("PSQL_CONNSTRING"), databaseConfig)
	err = database.Psql{}.Setup(databaseClient)
	if err != nil {
		app.logger.Fatalf("Error setting up new DatabaseClient: %v", err)
	}
	app.logger.Info("Database setup successful!")
	os.Exit(0)
}
