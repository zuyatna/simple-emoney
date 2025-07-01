package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-emoney/config"
	"simple-emoney/internal/app/handler"
	"simple-emoney/internal/app/middleware"
	"simple-emoney/internal/app/repository"
)

func SetupRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	transactionHandler *handler.TransactionHandler,
	redisRepo repository.RedisRepository,
) *gin.Engine {
	r := gin.Default()

	// public routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	r.POST("/register", authHandler.RegisterUser)
	r.POST("/login", authHandler.LoginUser)

	// authenticated routes
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(cfg, redisRepo))
	{
		users := v1.Group("/users")
		{
			users.POST("/topup", userHandler.TopUpBalance)
			users.GET("/balance", userHandler.GetUserBalance)
		}

		transactions := v1.Group("/transactions")
		{
			transactions.POST("/transfer", transactionHandler.Transfer)
			transactions.POST("/history", transactionHandler.GetTransactionHistory)
		}
	}

	return r
}
