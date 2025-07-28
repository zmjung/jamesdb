package disk

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jszwec/csvutil"
	"github.com/zmjung/jamesdb/graph"
)

func ReadNodesFromFile(filePath string) ([]graph.Node, error) {
	// This function reads nodes from a CSV file and returns them as a slice of graph.Node.

	fmt.Printf("Reading from file: %s\n", filePath)

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := string(file)

	csvReader := csv.NewReader(strings.NewReader(lines))

	dec, err := csvutil.NewDecoder(csvReader)
	dec.Tag = "json" // Use JSON tags for decoding
	dec.WithUnmarshalers(
		csvutil.NewUnmarshalers(
			csvutil.UnmarshalFunc(decodeList),
			csvutil.UnmarshalFunc(decodeMap),
		),
	)
	if err != nil {
		return nil, err
	}

	// TODO: is it more optimal to use slice or array?
	nodes := []graph.Node{}
	if err := dec.Decode(&nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func WriteNodeToFile(filePath string, node *graph.Node) error {
	// This function writes a CSV line to a file.
	// Assume that the filePath exists and is writable.

	fmt.Printf("Writing to file: %s\n", filePath)

	// Create a CSV writer
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Create an encoder
	encoder := csvutil.NewEncoder(writer)
	encoder.AutoHeader = false

	encoder.WithMarshalers(
		csvutil.NewMarshalers(
			csvutil.MarshalFunc(encodeList),
			csvutil.MarshalFunc(encodeMap),
		),
	)
	return encoder.Encode(node)
}

func encodeList(list []string) ([]byte, error) {
	if len(list) == 0 {
		return nil, nil
	}
	return []byte("[\"" + strings.Join(list, "\",\"") + "\"]"), nil
}

func decodeList(data []byte, list *[]string) error {
	if len(data) == 0 {
		return nil
	}

	str := strings.Trim(string(data), "[]")
	if str == "" {
		return nil
	}

	items := strings.Split(str, "\",\"")
	for _, item := range items {
		*list = append(*list, strings.Trim(item, "\""))
	}
	return nil
}

func encodeMap(kv map[string]string) ([]byte, error) {
	if len(kv) == 0 {
		return nil, nil
	}

	var sb strings.Builder
	sb.WriteString("{")
	count := len(kv)
	for k, v := range kv {
		sb.WriteString(fmt.Sprintf("\"%s\":\"%s\"", k, v))
		count--
		if count > 0 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("}")
	return []byte(sb.String()), nil
}

func decodeMap(data []byte, kv *map[string]string) error {
	if len(data) == 0 {
		return nil
	}

	str := strings.Trim(string(data), "{}")
	if str == "" {
		return nil
	}

	*kv = make(map[string]string)

	items := strings.Split(str, "\",\"")
	for _, item := range items {
		parts := strings.SplitN(item, "\":\"", 2)
		if len(parts) != 2 {
			continue // skip malformed entries
		}
		key := strings.Trim(parts[0], "\"")
		value := strings.Trim(parts[1], "\"")
		(*kv)[key] = value
	}
	return nil
}

func WriteCsvToFile(filePath string, csvLine string) error {
	// This function writes the CSV line to a file.
	// Assume that the filePath exists and is writable.

	fmt.Printf("Writing to file: %s\nCSV Line: %s\n", filePath, csvLine)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(csvLine); err != nil {
		return err
	}

	return err
}

func IsFileEmpty(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err == nil {
		// File exists, check size
		return info.Size() == 0, nil
	} else if os.IsNotExist(err) {
		// File does not exist, consider it empty and create it
		file, err := os.Create(filePath)
		if err != nil {
			return false, err
		}
		defer file.Close()
		return true, nil
	}
	return false, err
}

func AddFolder(rootPath string, folderName string) (string, error) {
	// Create full path
	absPath := filepath.Join(rootPath, folderName)

	// Create directory with full permissions
	// TODO: reevaluate permissions based on security needs
	if err := os.MkdirAll(absPath, os.ModePerm); err != nil {
		return "", err
	}

	return absPath, nil
}

func GetFilePath(rootPath string, fileName string) string {
	// This function returns the full file path for a given file name.
	// It combines the root path with the file name.
	return filepath.Join(rootPath, fileName)
}
