package grapher

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/disk"
)

const (
	TwoNodesCsv = `1,type1,node1,"[""edge1"",""edge2""]","{""trait1"":""value1""}"
2,type2,node2,"[""edge3"",""edge4""]","{""trait2"":""value2""}"
`
)

type nopWriteCloser struct {
	writer io.Writer
}

func (w *nopWriteCloser) Close() error {
	return nil
}

func (w *nopWriteCloser) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{
		writer: w,
	}
}

type MockFileAccessor struct {
	reader io.ReadCloser
	writer io.WriteCloser
	f      disk.FileAccessor
}

func GetFileAccessor(r io.Reader, w io.Writer) disk.FileAccessor {
	return getFileAccessorClosable(
		io.NopCloser(r),
		NopWriteCloser(w),
		disk.NewFileAccessor(),
	)
}

func getFileAccessorClosable(r io.ReadCloser, w io.WriteCloser, f disk.FileAccessor) disk.FileAccessor {
	return &MockFileAccessor{
		reader: r,
		writer: w,
		f:      f,
	}
}

func (m *MockFileAccessor) GetFileReader(filePath string) (io.ReadCloser, error) {
	return m.reader, nil
}

func (m *MockFileAccessor) GetFileWriter(filePath string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return m.writer, nil
}

func (m *MockFileAccessor) IsFileEmpty(filePath string) (bool, error) {
	return m.reader == nil, nil
}

func (m *MockFileAccessor) AddFolder(rootPath string, folderName string) (string, error) {
	return "", nil
}

func (m *MockFileAccessor) GetFilePath(rootPath string, fileName string) string {
	return m.f.GetFilePath(rootPath, fileName)
}

func TestWriteNodes(t *testing.T) {
	ctx := context.Background()

	reader := new(bytes.Buffer)
	writer := new(bytes.Buffer)
	f := GetFileAccessor(reader, writer)
	csv := disk.NewCsvAccessor(f)

	w := newWorker(f, csv, "nodePath", "nodeType")

	nodes := getTwoNodes()
	w.WriteNodes(ctx, nodes)

	bytes, err := io.ReadAll(writer)
	require.NoError(t, err)
	require.Equal(t, TwoNodesCsv, string(bytes))
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
