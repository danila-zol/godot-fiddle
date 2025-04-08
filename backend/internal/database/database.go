package database

import "gamehangar/internal/config"

// More of a guideline on how to setup a DatabaseClient
type DatabaseClient interface {
	NewDatabaseClient(connstring string, config *config.DatabaseConfig) any
	Setup() error
}
