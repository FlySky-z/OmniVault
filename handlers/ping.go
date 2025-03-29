package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler handles the ping request and responds with a "pong" message.
// @Summary PingHandler endpoint
// @Description Responds with a "pong" message to indicate the service is running.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
