package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/internal/handler"
)

type Router struct {
	GraphHandler *handler.GraphHandler
}

func NewRouter(gh *handler.GraphHandler) *Router {
	return &Router{
		GraphHandler: gh,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	healthRouter := engine.Group("/health")
	{
		healthRouter.GET("", handler.GetHealth)
	}

	graphRouter := engine.Group("/api/v1/graph")
	{
		graphRouter.GET("/node/:type", r.GraphHandler.GetGraphNodes)
		// readRouter.GET("/node/:nodeid", handler.GetGraphNode)

		graphRouter.POST("/node", r.GraphHandler.WriteGraphNode)
	}
}
