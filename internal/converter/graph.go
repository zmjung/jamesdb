package converter

import (
	"fmt"

	"github.com/zmjung/jamesdb/graph"
)

func ConvertToCSV(node *graph.Node) string {
	// This function converts a graph node to a CSV format string.
	// For now, we will return a placeholder string.
	// ID,Type,Name,Edges,Traits
	csvData := fmt.Sprintf("%s,%s,%s,\"%v\",\"%v\"\n", node.ID, node.Type, node.Name, node.Edges, node.Traits)
	return csvData
}
