package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-emoney/config"
)

func SetupRouter(
	cfg *config.Config,
) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	return r
}
