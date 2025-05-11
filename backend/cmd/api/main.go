package main

import (
	"context"
	"gamehangar/internal/config"
	"gamehangar/internal/config/psqlDatabseConfig"
	"gamehangar/internal/database/psqlDatabase"
	"gamehangar/internal/delivery/http/v1/handlers"
	"gamehangar/internal/delivery/http/v1/routes"
	"gamehangar/internal/enforcer/psqlCasbinClient"
	"gamehangar/internal/repository/psqlRepository"
	"gamehangar/internal/services"
	"gamehangar/pkg/ternMigrate"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
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
	validator *validator.Validate
	appRouter *echo.Router
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
	v := validator.New(validator.WithRequiredStructEnabled())
	getEnv()

	cfg := &appConfig{
		port: ":" + os.Getenv("PORT"),
	}

	app := &application{
		appConfig: cfg,
		echo:      e,
		validator: v,
		logger:    e.Logger,
	}

	// Might not look pretty to avoid passing around a pointer, but Psql*{} structs are empty
	databaseConfig, err := psqlDatabseConfig.PsqlConfig{}.NewConfig(
		psqlDatabase.MigrationFiles, os.Getenv("PSQL_MIGRATE_ROOT_DIR"),
	)
	if err != nil {
		app.logger.Fatalf("Error loading PSQL database Config: %v", err)
	}
	databaseClient, err := psqlDatabase.PsqlDatabase{}.NewDatabaseClient(
		os.Getenv("PSQL_CONNSTRING"), ternMigrate.Migrator{}, databaseConfig,
	)
	if err != nil {
		app.logger.Fatalf("Error setting up new DatabaseClient: %v", err)
	}
	app.logger.Info("Database setup successful!")

	wd, _ := os.Getwd()
	ce, err := psqlCasbinClient.CasbinConfig{}.NewCasbinClient(
		os.Getenv("PSQL_CONNSTRING"),
		wd+"/internal/enforcer/psqlCasbinClient/rbac_model.conf",
	)
	if err != nil {
		app.logger.Fatalf("Error setting up enforcer: %v", err)
	}

	userRepo := psqlRepository.NewPsqlUserRepository(databaseClient, ce)
	userAuthorizer := services.NewUserAuthorizer(userRepo, ce)
	userHandler := handlers.NewUserHandler(e, userRepo, app.validator, userAuthorizer)
	routes.NewUserRoutes(userHandler).InitRoutes(app.echo)

	assetRepo := psqlRepository.NewPsqlAssetRepository(databaseClient)
	assetHandler := handlers.NewAssetHandler(e, assetRepo, app.validator)
	routes.NewAssetRoutes(assetHandler).InitRoutes(app.echo)

	forumRepo := psqlRepository.NewPsqlForumRepository(databaseClient)
	forumHandler := handlers.NewForumHandler(e, forumRepo, app.validator)
	routes.NewForumRoutes(forumHandler).InitRoutes(app.echo)

	demoRepo := psqlRepository.NewPsqlDemoRepository(databaseClient)
	demoThreadSyncer := services.NewThreadSyncer(forumRepo, demoRepo, 1)
	demoHandler := handlers.NewDemoHandler(e, demoRepo, app.validator, demoThreadSyncer)
	routes.NewDemoRoutes(demoHandler).InitRoutes(app.echo)

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
