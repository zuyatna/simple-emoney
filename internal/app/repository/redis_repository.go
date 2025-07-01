package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"simple-emoney/internal/model"
	"time"
)

type RedisRepository interface {
	SetUserCache(user *model.User, duration time.Duration) error
	GetUserCache(userID string) (*model.User, error)
	DeleteUserCache(userID string) error
	SetAuthToken(token string, userID string, duration time.Duration) error
	GetAuthToken(token string) (string, error)
	DeleteAuthToken(token string) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r redisRepository) SetUserCache(user *model.User, duration time.Duration) error {
	key := fmt.Sprintf("user:%s", user.ID.String())
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user for cache: %w", err)
	}
	return r.client.Set(context.Background(), key, userData, duration).Err()
}

func (r redisRepository) GetUserCache(userID string) (*model.User, error) {
	key := fmt.Sprintf("user:%s", userID)
	userData, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // not found in cache
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user data from cache: %w", err)
	}

	var user model.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data from cache: %w", err)
	}
	return &user, nil
}

func (r redisRepository) DeleteUserCache(userID string) error {
	key := fmt.Sprintf("user:%s", userID)
	return r.client.Del(context.Background(), key).Err()
}

func (r redisRepository) SetAuthToken(token string, userID string, duration time.Duration) error {
	key := fmt.Sprintf("token:%s", token)
	return r.client.Set(context.Background(), key, userID, duration).Err()
}

func (r redisRepository) GetAuthToken(token string) (string, error) {
	key := fmt.Sprintf("token:%s", token)
	userID, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get auth token from Redis: %w", err)
	}
	return userID, nil
}

func (r redisRepository) DeleteAuthToken(token string) error {
	key := fmt.Sprintf("token:%s", token)
	return r.client.Del(context.Background(), key).Err()
}
