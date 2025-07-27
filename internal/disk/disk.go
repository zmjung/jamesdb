package disk

import (
	"fmt"
	"os"
	"path/filepath"
)

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
