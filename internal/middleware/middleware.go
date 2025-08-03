package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/internal/uuid"
)

func GetLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestId := setRequestId(c)
		log.Printf("[%s] %s %s started", requestId, c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("[%s] %s %s completed in %v", requestId, c.Request.Method, c.Request.URL.Path, time.Since(start))
	}
}

func GetRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(500, gin.H{"message": "Internal server error"})
			}
		}()
		c.Next()
	}
}

func setRequestId(c *gin.Context) string {
	requestId, err := uuid.GenerateShortID()
	if err != nil {
		log.Printf("Error generating thread ID: %v", err)
		requestId = "unknown"
	}
	c.Set("requestId", requestId)
	return requestId
}
