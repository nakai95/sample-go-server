package main

import (
	"fmt"
	"log"
	"os"
	"sample-go-server/constants"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv(constants.DB_USER)
	dbPassword := os.Getenv(constants.DB_PASSWORD)
	dbName := os.Getenv(constants.DB_NAME)
	dbHost := os.Getenv(constants.DB_HOST)
	dbPort := os.Getenv(constants.DB_PORT)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	m, err := migrate.New(
		"file://migrate/postgres/migrations",
		dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
