package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"simple-emoney/config"
	"time"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		_, err := rdb.Ping(ctx).Result()
		if err == nil {
			log.Println("Successfully connected to redis!")
			return rdb, nil
		}

		log.Printf("Waiting for redis connection, attempt %d: %v\n\n", i+1, err)
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to redis after multiple retries")
}
