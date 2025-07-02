package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"simple-emoney/internal/model"
	"time"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id uuid.UUID) (*model.User, error)
	UpdateUserBalance(tx *sql.Tx, userID uuid.UUID, amount float64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (u userRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (username, password_hash, balance) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := u.db.QueryRow(query, user.Username, user.PasswordHash, user.Balance).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (u userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, password_hash, balance, created_at, updated_at FROM users WHERE username = $1`
	err := u.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user by username: %w", err)
	}
	return &user, nil
}

func (u userRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, password_hash, balance, created_at, updated_at FROM users WHERE id = $1`
	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}
	return &user, nil
}

func (u userRepository) UpdateUserBalance(tx *sql.Tx, userID uuid.UUID, amount float64) error {
	log.Printf("Executing UpdateUserBalance for user %s with amount %.2f", userID.String(), amount)
	query := `UPDATE users SET balance = balance + $1, updated_at = $2 WHERE id = $3`
	_, err := tx.Exec(query, amount, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("error updating user balance: %w", err)
	}
	return nil
}
