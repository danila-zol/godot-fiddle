package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	database "game-hangar/database"
	"game-hangar/server"
)

func getDsn() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	return dsn
}

func main() {
	database.SetupDB(getDsn())
	server.Setup()
}
