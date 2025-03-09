package main

import (
	database "auction-api/database/migrations"
	"database/sql"
	"log"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal(err)
	}
}
