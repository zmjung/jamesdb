package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/internal/handler"
)

func SetupRoutes(r *gin.Engine) {
	healthRouter := r.Group("/health")
	{
		healthRouter.GET("", handler.GetHealth)
	}

	readRouter := r.Group("/api/v1/graph")
	{
		readRouter.GET("/node", handler.GetGraphNodes)
		// readRouter.GET("/node/:nodeid", handler.GetGraphNode)
	}
}
