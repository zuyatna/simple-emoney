package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simple-emoney/internal/app/service"
	"simple-emoney/internal/model"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (th *TransactionHandler) Transfer(c *gin.Context) {
	var req model.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	senderID := c.GetString("userID")
	if senderID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sender ID not found in context"})
		return
	}

	err := th.transactionService.Transfer(senderID, &req)
	if err != nil {
		log.Printf("Error during transfer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer succeeded",
	})
}

func (th *TransactionHandler) GetTransactionHistory(c *gin.Context) {
	userID := c.GetString("userID") // get user ID from authenticated context
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not found in context"})
		return
	}

	transactions, err := th.transactionService.GetTransactionHistory(userID)
	if err != nil {
		log.Printf("Error getting transaction history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}
