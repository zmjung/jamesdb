package disk

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zmjung/jamesdb/graph"
)

const csvTwoNodes = `id,type,name,edges,traits
1,type1,node1,"[""edge1"",""edge2""]","{""trait1"":""value1""}"
2,type2,node2,"[""edge3"",""edge4""]","{""trait2"":""value2""}"
`

func Test_readNodesFromIo(t *testing.T) {
	csv := &csvService{}
	ctx := context.Background()

	reader := strings.NewReader(csvTwoNodes)
	nodes, err := csv.readNodesFromReader(ctx, reader)
	if err != nil {
		t.Fatalf("Failed to read nodes from IO: %v", err)
	}

	expected := getTwoNodes()
	require.Equal(t, expected, nodes)

	// Additional check: Ensure that the number of nodes matches the expected count
	if len(nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %d. Possible CSV or JSON decoding error.", len(nodes))
	}
}

func Test_writeCsvToWriter(t *testing.T) {
	csv := &csvService{}
	ctx := context.Background()

	nodes := getTwoNodes()

	writer := new(strings.Builder)
	writer.WriteString(graph.NodeCsvHeader)

	err := csv.writeCsvToWriter(ctx, writer, nodes)
	if err != nil {
		t.Fatalf("Failed to write nodes to writer: %v", err)
	}

	expected := csvTwoNodes
	actual := writer.String()
	require.Equal(t, expected, actual)
	// Additional check: Ensure that the output is not empty
	if actual == "" {
		t.Fatalf("Expected non-empty output, got empty string. Possible CSV encoding error.")
	}
}

func getTwoNodes() []graph.Node {
	return []graph.Node{
		{
			ID:     "1",
			Type:   "type1",
			Name:   "node1",
			Edges:  []string{"edge1", "edge2"},
			Traits: map[string]string{"trait1": "value1"},
		},
		{
			ID:     "2",
			Type:   "type2",
			Name:   "node2",
			Edges:  []string{"edge3", "edge4"},
			Traits: map[string]string{"trait2": "value2"},
		},
	}
}
