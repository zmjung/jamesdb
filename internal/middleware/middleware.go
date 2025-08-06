package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/internal/log"
)

func GetLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := log.ConvertContext(c)
		slog.InfoContext(ctx, "Started request")
		c.Next()
		slog.InfoContext(ctx, "Completed request", "time", time.Since(start))
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
