package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/graph"
)

type GraphHandler struct {
	StorageRootPath string
}

func NewGraphHandler(cfg config.Config) *GraphHandler {
	return &GraphHandler{
		StorageRootPath: cfg.Database.RootPath,
	}
}

func (gh *GraphHandler) GetGraphNodes(c *gin.Context) {
	// TODO: sanitize clusterId input
	clusterId := c.Param("clusterId")
	fmt.Printf("%s\n", gh.StorageRootPath+"/nodes/"+clusterId+".csv")
	// This is a placeholder for the actual logic to retrieve graph nodes.
	// For now, we will return a static list of nodes.
	nodes := []graph.Node{
		{
			ID:     "1",
			Name:   "Node 1",
			Edges:  []string{"2", "3"},
			Traits: map[string]string{"color": "red", "size": "large"},
		},
	}
	c.JSON(200, nodes)
}
