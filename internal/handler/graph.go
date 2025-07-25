package handler

import (
	"github.com/gin-gonic/gin"
	graph "github.com/zmjung/jamesdb/graph"
)

func GetGraphNodes(c *gin.Context) {
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
