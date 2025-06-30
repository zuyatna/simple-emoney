package repository

import (
	"github.com/redis/go-redis/v9"
	"simple-emoney/internal/model"
	"time"
)

type RedisRepository interface {
	SetUserCache(user *model.User, duration time.Duration) error
	GetUserCache(userID string) (*model.User, error)
	DeleteUserCache(userID string) error
	SetAuthToken(token string, userID string, duration time.Duration) error
	GetUserAuthToken(token string) (string, error)
	DeleteAuthToken(token string) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r redisRepository) SetUserCache(user *model.User, duration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r redisRepository) GetUserCache(userID string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r redisRepository) DeleteUserCache(userID string) error {
	//TODO implement me
	panic("implement me")
}

func (r redisRepository) SetAuthToken(token string, userID string, duration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r redisRepository) GetUserAuthToken(token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r redisRepository) DeleteAuthToken(token string) error {
	//TODO implement me
	panic("implement me")
}
