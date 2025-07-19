package graph

type Node struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Edges []string `json:"edges"`
	Traits map[string]string `json:"traits"`
}
