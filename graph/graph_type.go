package graph

type CsvHeader string

// do not create this dynamically
// need this to be this exact order for data sanity / backwards compatibility
// I already know about `csvutil.Header(Node{}, "json")`
const NodeCsvHeader = "id,type,name,edges,traits\n"

type Node struct {
	ID     string            `json:"id"`
	Type   string            `json:"type" binding:"required"`
	Name   string            `json:"name" binding:"required"`
	Edges  []string          `json:"edges,omitempty"`
	Traits map[string]string `json:"traits,omitempty"`
}
