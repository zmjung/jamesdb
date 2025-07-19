package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string {"status": "healthy"})
}