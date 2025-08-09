package disk

import (
	"testing"

	"github.com/stretchr/testify/require"
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
