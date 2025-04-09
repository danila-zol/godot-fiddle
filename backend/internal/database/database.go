package database

import "gamehangar/internal/config"

type DatabaseClientCreator interface {
	NewDatabaseClient(connstring string, config *config.DatabaseConfig) any
	Setup() error
}
