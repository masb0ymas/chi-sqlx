package main

import (
	"chi-sqlx/config"
	"chi-sqlx/database"
	"chi-sqlx/routes"
	"log"
)

func main() {
	dbname := config.Env("DB_DATABASE", "db_example")

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()
	log.Printf("successfully connected to database %v", dbname)

	// register routes
	routes.RegisterRoutes(db.GetDB())
}
