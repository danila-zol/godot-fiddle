package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	database "game-hangar/database"
	_ "game-hangar/docs"
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

//	@title			Game Hangar
//	@version		1.0
//	@description	A backend for asset catalogue
//	@contact.name	Mikhail Pecherkin
//	@contact.email	m.pecherkin.sas@gmail.com
//	@host			localhost:9938
//	@BasePath		/

func main() {
	database.SetupDB(getDsn())
	server.Setup(os.Getenv("HOST"))
}
