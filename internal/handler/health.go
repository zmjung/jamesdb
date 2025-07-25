package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
