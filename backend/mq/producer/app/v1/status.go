package v1

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /v1

// Ping
// @Just a ping
// @Schemes
// @Description do ping
// @Tags health checks
// @Accept json
// @Produce json
// @Success 200 {json} message:pong
// @Router /ping [get]
// @Summary Liveness check
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
