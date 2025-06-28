package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"simple-emoney/config"
	"time"
)

func NewPostgreSQLDB(cfg *config.Config) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to database")
			return db, nil
		}

		log.Printf("Waiting for connection to database to be ready, attempt %d: %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return db, fmt.Errorf("error connecting to database: %v", err)
}
