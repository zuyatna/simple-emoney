package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"simple-emoney/internal/app/repository"
	"simple-emoney/internal/model"
	"time"
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
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		log.Println("Recovered from panic during top-up, rolling back transaction")
		if err != nil {
			return
		}
	}(tx) // rollback on error

	err = s.userRepo.UpdateUserBalance(tx, userID, req.Amount)
	if err != nil {
		return fmt.Errorf("failed to top-up user balance: %w", err)
	}

	transaction := &model.Transaction{
		SenderID:        userID,
		ReceiverID:      userID,
		Amount:          req.Amount,
		TransactionType: "topup",
	}

	err = s.transactionRepo.CreateTransaction(tx, transaction)
	if err != nil {
		return fmt.Errorf("failed to record top-up transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit top-up transaction: %w", err)
	}

	// invalidate cache after update
	err = s.redisRepo.DeleteUserCache(userID.String())
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserBalance(userID string) (float64, error) {
	// try to get from cache first
	cachedUser, err := s.redisRepo.GetUserCache(userID)
	if err != nil {
		log.Printf("Error getting user from Redis cache: %v", err)
	}
	if cachedUser != nil {
		return cachedUser.Balance, nil
	}

	// if not in cache, get from DB
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetUserByID(userUUID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("user not found")
	}

	// cache the user for future requests
	err = s.redisRepo.SetUserCache(user, 1*time.Hour)
	if err != nil {
		return 0, err
	}

	return user.Balance, nil
}
