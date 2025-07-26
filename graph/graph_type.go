package graph

type Node struct {
	ID     string            `json:"id"`
	Type   string            `json:"Type" binding:"required"`
	Name   string            `json:"name" binding:"required"`
	Edges  []string          `json:"edges"`
	Traits map[string]string `json:"traits"`
}
