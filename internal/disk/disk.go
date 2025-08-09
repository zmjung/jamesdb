package disk

import (
	"io"
	"os"
	"path/filepath"
)

type FileAccessor interface {
	GetFileReader(filePath string) (io.ReadCloser, error)
	GetFileWriter(filePath string, flag int, perm os.FileMode) (io.WriteCloser, error)
	IsFileEmpty(filePath string) (bool, error)
	AddFolder(rootPath string, folderName string) (string, error)
	GetFilePath(rootPath string, fileName string) string
}

type filer struct{}

func NewFileAccessor() FileAccessor {
	return &filer{}
}

func (f *filer) GetFileReader(filePath string) (io.ReadCloser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *filer) GetFileWriter(filePath string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	file, err := os.OpenFile(filePath, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *filer) IsFileEmpty(filePath string) (bool, error) {
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

func (f *filer) AddFolder(rootPath string, folderName string) (string, error) {
	// Create full path
	absPath := filepath.Join(rootPath, folderName)

	// Create directory with full permissions
	// TODO: reevaluate permissions based on security needs
	if err := os.MkdirAll(absPath, os.ModePerm); err != nil {
		return "", err
	}

	return absPath, nil
}

func (f *filer) GetFilePath(rootPath string, fileName string) string {
	// This function returns the full file path for a given file name.
	// It combines the root path with the file name.
	return filepath.Join(rootPath, fileName)
}
