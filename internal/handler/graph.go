package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/grapher"
	"github.com/zmjung/jamesdb/internal/log"
	"github.com/zmjung/jamesdb/internal/uuid"
)

type GraphHandler struct {
	StorageRootPath string
	Grapher         grapher.Grapher
}

func NewGraphHandler(cfg *config.Config, g grapher.Grapher) *GraphHandler {
	return &GraphHandler{
		StorageRootPath: cfg.Database.RootPath,
		Grapher:         g,
	}
}

func (gh *GraphHandler) GetGraphNodes(c *gin.Context) {
	// This function gets all graph nodes according to the passed node type
	ctx := log.ConvertContext(c)
	// TODO: sanitize type input
	nodeType := c.Param("type")
	nodes, err := gh.Grapher.ReadNodesByType(ctx, nodeType)

	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to retrieve nodes of type %s: %v", nodeType, err)})
		return
	}
	c.JSON(200, nodes)
}

func (gh *GraphHandler) CreateGraphNode(c *gin.Context) {
	// This function writes a graph node to the storage.
	ctx := log.ConvertContext(c)

	node := &graph.Node{}
	if err := c.ShouldBindJSON(node); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input", "node": node})
		return
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err)
		c.JSON(500, gin.H{"error": "Failed to generate UUID", "node": node})
		return
	}
	node.ID = id

	err = gh.Grapher.WriteNode(ctx, node)
	if err != nil {
		fmt.Printf("Error writing node data: %v\n", err)
		c.JSON(500, gin.H{"error": "Failed to write node data", "node": node})
		return
	}

	c.JSON(200, gin.H{"message": "Node data written successfully", "node": node})
}
