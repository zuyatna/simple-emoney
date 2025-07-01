package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simple-emoney/internal/app/service"
	"simple-emoney/internal/model"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (uh *UserHandler) TopUpBalance(c *gin.Context) {
	var req model.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// in real app, userID should come from authenticated context, not body
	// for simplicity, using req.UserID for now.
	// you might want to get it from c.GetString("userID") after auth middleware.

	err := uh.userService.TopUpBalance(&req)
	if err != nil {
		log.Printf("Error topping up balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Balance topped up successfully",
	})
}
