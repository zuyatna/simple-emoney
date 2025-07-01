package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"simple-emoney/config"
	"simple-emoney/internal/app/repository"
	"simple-emoney/pkg/utils"
	"strings"
)

func AuthMiddleware(cfg *config.Config, redisRepo repository.RedisRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is empty",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is invalid",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// optional: check if token is blacklisted in Redis (e.g., after logout)
		_, err := redisRepo.GetAuthToken(tokenString)
		if err != nil && !errors.Is(err, redis.Nil) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check token validity",
			})
			c.Abort()
			return
		}
		if err == nil {
			// token found in Redis (meaning it's valid for simple storage, or blacklisted for invalidate)
			// for blacklist, if token is found here, it means it's blacklisted
			// if we use Redis to store active tokens, then finding it means it's active.
			// let's assume for this example, if token exists in redis it means it's active/valid.
			// more robust blacklist requires checking if token is explicitly marked as blacklisted.
		}

		claims, err := utils.VerifyJWTToken(tokenString, cfg.JWTSecretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// set user ID and username in context for handlers to use
		c.Set("userID", claims.UserID.String())
		c.Set("username", claims.Username)
		c.Next()
	}
}
