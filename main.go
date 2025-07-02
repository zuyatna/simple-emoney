package main

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"log"
	"simple-emoney/config"
	"simple-emoney/internal/app/handler"
	"simple-emoney/internal/app/repository"
	"simple-emoney/internal/app/service"
	"simple-emoney/internal/router"
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

	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			log.Fatalf("Failed to close redis: %v", err)
		}
	}(redisClient)

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)

	authService := service.NewAuthService(userRepo, redisRepo, cfg)
	userService := service.NewUserService(db, userRepo, transactionRepo, redisRepo)
	transactionService := service.NewTransactionService(db, userRepo, transactionRepo, redisRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	r := router.SetupRouter(cfg, authHandler, userHandler, transactionHandler, redisRepo)

	log.Printf("Server is running on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
