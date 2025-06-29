package model

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID              uuid.UUID `json:"id"`
	SenderID        uuid.UUID `json:"sender_id"`
	ReceiverID      uuid.UUID `json:"receiver_id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}

type TransferRequest struct {
	ReceiverUsername string  `json:"receiver_username" binding:"required"`
	Amount           float64 `json:"amount" binding:"required,gt=0"`
}
