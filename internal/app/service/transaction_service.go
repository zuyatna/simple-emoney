package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"simple-emoney/internal/app/repository"
	"simple-emoney/internal/model"
)

type TransactionService interface {
	Transfer(senderID string, req *model.TransferRequest) error
	GetTransactionHistory(userID string) ([]model.Transaction, error)
}

type transactionService struct {
	db              *sql.DB // inject DB for manual transaction management
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
	redisRepo       repository.RedisRepository
}

func NewTransactionService(db *sql.DB, userRepo repository.UserRepository, transactionRepo repository.TransactionRepository, redisRepo repository.RedisRepository) TransactionService {
	return &transactionService{
		db:              db,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		redisRepo:       redisRepo,
	}
}

func (t transactionService) Transfer(senderID string, req *model.TransferRequest) error {
	senderIDStr, err := uuid.Parse(senderID)
	if err != nil {
		return errors.New("invalid sender ID format")
	}

	// start a database transaction for atomicity
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for transfer: %w", err)
	}
	// NEW: flag to track if the transaction has been committed
	committed := false
	defer func() {
		if !committed { // only rollback if not already committed
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Failed to rollback transaction: %v", rbErr) // log the rollback error, don't return it
			}
		}
	}()

	// get sender (for update, use FOR UPDATE in real scenario to prevent race condition)
	sender, err := t.userRepo.GetUserByID(senderIDStr) // in a real scenario, fetch with FOR UPDATE
	if err != nil {
		return err
	}
	if sender == nil {
		return errors.New("sender not found")
	}

	// get receiver by username
	receiver, err := t.userRepo.GetUserByUsername(req.ReceiverUsername)
	if err != nil {
		return err
	}
	if receiver == nil {
		return errors.New("receiver not found")
	}

	if sender.ID == receiver.ID {
		return errors.New("cannot transfer to self")
	}

	if sender.Balance < req.Amount {
		return errors.New("insufficient balance")
	}

	log.Printf("Transferring amount %.2f from sender %s to receiver %s", req.Amount, sender.Username, receiver.Username)

	// deduct from sender's balance
	err = t.userRepo.UpdateUserBalance(tx, sender.ID, -req.Amount)
	if err != nil {
		return fmt.Errorf("failed to deduct from sender balance: %w", err)
	}

	// add to receiver's balance
	err = t.userRepo.UpdateUserBalance(tx, receiver.ID, req.Amount)
	if err != nil {
		// log specific error and ensure rollback is handled by defer
		return fmt.Errorf("failed to add to receiver balance: %w", err)
	}

	transaction := &model.Transaction{
		SenderID:        sender.ID,
		ReceiverID:      receiver.ID,
		Amount:          req.Amount,
		TransactionType: "transfer",
	}
	err = t.transactionRepo.CreateTransaction(tx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// invalidate caches for both sender and receiver
	err = t.redisRepo.DeleteUserCache(sender.ID.String())
	if err != nil {
		return err
	}

	err = t.redisRepo.DeleteUserCache(receiver.ID.String())
	if err != nil {
		return err
	}

	return nil
}

func (t transactionService) GetTransactionHistory(userID string) ([]model.Transaction, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := t.userRepo.GetUserByID(uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// for transaction history, it might not be beneficial to cache directly as it changes ofter
	transactions, err := t.transactionRepo.GetTransactionByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions history: %w", err)
	}
	return transactions, nil
}
