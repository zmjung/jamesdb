package graph

type Node struct {
	ID     string            `json:"id"`
	Type   string            `json:"type" binding:"required"`
	Name   string            `json:"name" binding:"required"`
	Edges  []string          `json:"edges,omitempty"`
	Traits map[string]string `json:"traits,omitempty"`
}
