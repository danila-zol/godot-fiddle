package main

import (
	"context"
	"gamehangar/internal/config"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/delivery/http/v1/handlers"
	"gamehangar/internal/delivery/http/v1/routes"
	"gamehangar/internal/repository/psqlRepository"
	"gamehangar/internal/services"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// Exclusive to package main since it is a central app config
type appConfig struct {
	port string
}

type application struct {
	appConfig *appConfig
	echo      *echo.Echo
	logger    echo.Logger
	appRouter *echo.Router
	validator echo.Validator
}

type DatabaseConfigCreator interface {
	NewConfig() (*config.DatabaseConfig, error)
}

type databaseClientCreator interface {
	NewDatabaseClient(connstring string, config *config.DatabaseConfig) (any, error)
}

func getEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

// @title			Game Hangar
// @version		1.0
// @description	A backend for Game Hangar game prototyping web service
// @contact.name	Mikhail Pecherkin
// @contact.email	m.pecherkin.sas@gmail.com
// @BasePath		/game-hangar
// @securityDefinitions.apikey ApiSessionCookie
// @in header
// @name sessionID
func main() {
	e := echo.New()
	getEnv()

	cfg := &appConfig{
		port: ":" + os.Getenv("PORT"),
	}

	app := &application{
		appConfig: cfg,
		echo:      e,
		logger:    e.Logger,
		validator: e.Validator,
	}

	// Might not look pretty to avoid passing around a pointer, but Psql*{} structs are empty
	databaseConfig, err := psqlDatabseConfig.PsqlConfig{}.NewConfig(
		psqlDatabase.MigrationFiles, os.Getenv("PSQL_MIGRATE_ROOT_DIR"),
	)
	if err != nil {
		app.logger.Fatalf("Error loading PSQL database Config: %v", err)
	}
	databaseClient, err := psqlDatabase.PsqlDatabase{}.NewDatabaseClient(
		os.Getenv("PSQL_CONNSTRING"), databaseConfig,
	)
	if err != nil {
		app.logger.Fatalf("Error setting up new DatabaseClient: %v", err)
	}
	app.logger.Info("Database setup successful!")

	assetRepo := psqlRepository.NewPsqlAssetRepository(databaseClient)
	assetHandler := handlers.NewAssetHandler(e, assetRepo)
	routes.NewAssetRoutes(assetHandler).InitRoutes(app.echo)

	forumRepo := psqlRepository.NewPsqlForumRepository(databaseClient)
	forumHandler := handlers.NewForumHandler(e, forumRepo)
	routes.NewForumRoutes(forumHandler).InitRoutes(app.echo)

	demoRepo := psqlRepository.NewPsqlDemoRepository(databaseClient)
	demoThreadSyncer := services.NewThreadSyncer(forumRepo, demoRepo, 1)
	demoHandler := handlers.NewDemoHandler(e, demoRepo, demoThreadSyncer)
	routes.NewDemoRoutes(demoHandler).InitRoutes(app.echo)

	userRepo := psqlRepository.NewPsqlUserRepository(databaseClient)
	userIDLookup := services.NewUserIdentifier(userRepo)
	userHandler := handlers.NewUserHandler(e, userRepo, userIDLookup)
	routes.NewUserRoutes(userHandler).InitRoutes(app.echo)

	app.appRouter = app.routes(app.echo)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := app.echo.Start(app.appConfig.port); err != nil && err != http.ErrServerClosed {
			app.logger.Fatal("Shutting down the server")
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		app.logger.Fatal(err)
	}
}
