package repository

import (
	"database/sql"
	"fmt"
	"simple-emoney/internal/model"
)

type TransactionRepository interface {
	CreateTransaction(tx *sql.Tx, transaction *model.Transaction) error
	GetTransactionByUserID(userID string) ([]model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (t transactionRepository) CreateTransaction(sqlTx *sql.Tx, transaction *model.Transaction) error {
	query := `INSERT INTO transactions (sender_id, receiver_id, amount, transaction_type) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := sqlTx.QueryRow(query, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.TransactionType).Scan(&transaction.ID, &transaction.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating transaction within transaction: %w", err)
	}
	return nil
}

func (t transactionRepository) GetTransactionByUserID(userID string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	query := `SELECT * FROM transactions WHERE sender_id = $1 OR receiver_id = $1 ORDER BY created_at DESC LIMIT 100`

	rows, err := t.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions by user ID: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.TransactionType, &transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction row: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return transactions, nil
}
