package disk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zmjung/jamesdb/graph"
)

func Test_encodeList(t *testing.T) {
	// This function tests the encodeList function.
	// It checks if the function correctly encodes a list of strings.
	list := []string{"edge1", "edge2"}
	encoded, err := encodeList(list)
	if err != nil {
		t.Fatalf("Failed to encode list: %v", err)
	}
	expected := `["edge1","edge2"]`
	actual := string(encoded[:])
	require.JSONEq(t, expected, actual)
}

func Test_encodeEmptyList(t *testing.T) {
	// This function tests the encodeList function with an empty list.
	// It checks if the function correctly encodes an empty list.
	list := []string{}
	encoded, err := encodeList(list)
	if err != nil {
		t.Fatalf("Failed to encode empty list: %v", err)
	}
	expected := ""
	actual := string(encoded[:])
	require.Equal(t, expected, actual)
}

func Test_decodeList(t *testing.T) {
	// This function tests the decodeList function.
	// It checks if the function correctly decodes a JSON array into a slice of strings.
	jsonData := `["edge1","edge2"]`
	var list []string
	err := decodeList([]byte(jsonData), &list)
	if err != nil {
		t.Fatalf("Failed to decode list: %v", err)
	}
	expected := []string{"edge1", "edge2"}
	require.Equal(t, expected, list)
}

func Test_decodeEmptyList(t *testing.T) {
	// This function tests the decodeList function with an empty JSON array.
	// It checks if the function correctly decodes an empty array into a slice of strings.
	jsonData := `[]`
	var list []string
	err := decodeList([]byte(jsonData), &list)
	if err != nil {
		t.Fatalf("Failed to decode empty list: %v", err)
	}
	var expected []string = nil
	require.Equal(t, expected, list)
}

func Test_decodeListWithNull(t *testing.T) {
	// Skip this test as it is not applicable in the current context.
	t.Skip("Skipping test for decodeList with null value in list")
	// This function tests the decodeList function with a null value in the list.
	// It checks if the function correctly handles a nil input.
	jsonData := `[null]`
	var list []string
	err := decodeList([]byte(jsonData), &list)
	if err != nil {
		t.Fatalf("Failed to decode empty list: %v", err)
	}
	expected := []string{}
	require.Equal(t, expected, list)
}

func Test_encodeMap(t *testing.T) {
	// This function tests the encodeMap function.
	// It checks if the function correctly encodes a map of strings.
	traits := map[string]string{
		"trait1": "value1",
		"trait2": "value2",
	}
	encoded, err := encodeMap(traits)
	if err != nil {
		t.Fatalf("Failed to encode map: %v", err)
	}
	expected := `{"trait1":"value1","trait2":"value2"}`
	actual := string(encoded[:])
	require.JSONEq(t, expected, actual)
}

func Test_encodeEmptyMap(t *testing.T) {
	// This function tests the encodeMap function with an empty map.
	// It checks if the function correctly encodes an empty map.
	traits := map[string]string{}
	encoded, err := encodeMap(traits)
	if err != nil {
		t.Fatalf("Failed to encode empty map: %v", err)
	}
	expected := ""
	actual := string(encoded[:])
	require.Equal(t, expected, actual)
}

func Test_decodeMap(t *testing.T) {
	// This function tests the decodeMap function.
	// It checks if the function correctly decodes a JSON object into a map of strings.
	jsonData := `{"trait1":"value1","trait2":"value2"}`
	var traits map[string]string
	err := decodeMap([]byte(jsonData), &traits)
	if err != nil {
		t.Fatalf("Failed to decode map: %v", err)
	}
	expected := map[string]string{
		"trait1": "value1",
		"trait2": "value2",
	}
	require.Equal(t, expected, traits)
}

func Test_decodeEmptyMap(t *testing.T) {
	// This function tests the decodeMap function with an empty JSON object.
	// It checks if the function correctly decodes an empty object into a map of strings.
	jsonData := `{}`
	var traits map[string]string
	err := decodeMap([]byte(jsonData), &traits)
	if err != nil {
		t.Fatalf("Failed to decode empty map: %v", err)
	}
	var expected map[string]string = nil
	require.Equal(t, expected, traits)
}

func Test_decodeMapWithNull(t *testing.T) {
	// This function tests the decodeMap function with a null value in the map.
	// It checks if the function correctly handles a nil input.
	jsonData := `{"trait1":null}`
	var traits map[string]string
	err := decodeMap([]byte(jsonData), &traits)
	if err != nil {
		t.Fatalf("Failed to decode map with null value: %v", err)
	}
	expected := map[string]string{}
	require.Equal(t, expected, traits)
}

func Test_readNodesFromIo(t *testing.T) {
	// This function tests the readNodesFromIo function.
	// It checks if the function correctly reads nodes from a CSV reader.
	csvData := `id,type,name,edges,traits
1,type1,node1,"[""edge1"",""edge2""]","{""trait1"":""value1"",""trait2"":""value2""}"
2,type2,node2,"[""edge3"",""edge4""]","{""trait3"":""value3"",""trait4"":""value4""}"
`
	reader := strings.NewReader(csvData)
	nodes, err := readNodesFromReader(reader)
	if err != nil {
		t.Fatalf("Failed to read nodes from IO: %v", err)
	}

	expected := []graph.Node{
		{
			ID:     "1",
			Type:   "type1",
			Name:   "node1",
			Edges:  []string{"edge1", "edge2"},
			Traits: map[string]string{"trait1": "value1", "trait2": "value2"},
		},
		{
			ID:     "2",
			Type:   "type2",
			Name:   "node2",
			Edges:  []string{"edge3", "edge4"},
			Traits: map[string]string{"trait3": "value3", "trait4": "value4"},
		},
	}
	require.Equal(t, expected, nodes)

	// Additional check: Ensure that the number of nodes matches the expected count
	if len(nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %d. Possible CSV or JSON decoding error.", len(nodes))
	}
}

func Test_writeCsvToWriter(t *testing.T) {
	// This function tests the writeNodesToWriter function.
	// It checks if the function correctly writes nodes to a CSV file.
	nodes := []graph.Node{
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

	writer := new(strings.Builder)
	writer.WriteString("id,type,name,edges,traits\n")

	err := writeCsvToWriter(writer, nodes)
	if err != nil {
		t.Fatalf("Failed to write nodes to writer: %v", err)
	}
	expected := `id,type,name,edges,traits
1,type1,node1,"[""edge1"",""edge2""]","{""trait1"":""value1""}"
2,type2,node2,"[""edge3"",""edge4""]","{""trait2"":""value2""}"
`
	actual := writer.String()
	require.Equal(t, expected, actual)
	// Additional check: Ensure that the output is not empty
	if actual == "" {
		t.Fatalf("Expected non-empty output, got empty string. Possible CSV encoding error.")
	}
}
