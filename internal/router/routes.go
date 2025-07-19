package router

import (
	"github.com/zmjung/jamesdb/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(e *gin.Engine) {
	// Health check endpoint
	e.GET("/health", handler.GetHealthCheck)
}