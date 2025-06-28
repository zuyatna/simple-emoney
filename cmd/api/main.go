package main

import (
	"database/sql"
	"log"
	"net/http"
	"simple-emoney/config"
	"simple-emoney/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgreSQLDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}(db)

	log.Printf("Server is running on port %s", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
