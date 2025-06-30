package service

import (
	"database/sql"
	"simple-emoney/internal/app/repository"
	"simple-emoney/internal/model"
)

type UserService interface {
	TopUpBalance(req *model.TopUpRequest) error
	GetUserBalance(userID string) (float64, error)
}

type userService struct {
	db              *sql.DB // inject DB for manual transaction management
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
	redisRepo       repository.RedisRepository
}

func NewUserService(db *sql.DB, userRepo repository.UserRepository, transactionRepo repository.TransactionRepository, redisRepo repository.RedisRepository) UserService {
	return &userService{
		db:              db,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		redisRepo:       redisRepo,
	}
}

func (s *userService) TopUpBalance(req *model.TopUpRequest) error {
	// TODO: Implement the logic to top up user balance
	return nil
}

func (s *userService) GetUserBalance(userID string) (float64, error) {
	// TODO: Implement the logic to get user balance
	return 0, nil
}
