package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/assaidy/goblog/repo/postgres_repo"
	"github.com/assaidy/goblog/router"
	"github.com/assaidy/goblog/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dbConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	repo, err := postgres_repo.NewPostgresRepo(dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := postgres_repo.Migrate(repo.DB); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
    defer repo.DB.Close()

	router := router.NewRouter(repo)

	log.Println("Running server on port %s", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))
}
