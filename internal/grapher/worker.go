package grapher

import (
	"context"
	"log/slog"
	"sync"

	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/disk"
)

var EmptyGraphNodes = []graph.Node{}

type Worker interface {
	ReadNodes(ctx context.Context) ([]graph.Node, error)
	WriteNodes(ctx context.Context, nodes []graph.Node) error
}

type worker struct {
	f        disk.FileAccessor
	csv      disk.CsvAccessor
	nodeType string
	filePath string
	lock     *sync.Mutex
}

func newWorker(f disk.FileAccessor, csv disk.CsvAccessor, nodePath string, nodeType string) Worker {
	filePath := f.GetFilePath(nodePath, nodeType+".csv")
	err := csv.CreateFileWithHeader(context.Background(), filePath, graph.NodeCsvHeader)
	if err != nil {
		slog.Error("Error creating headers for nodes file", "error", err)
		return nil
	}

	return &worker{
		f:        f,
		csv:      csv,
		nodeType: nodeType,
		filePath: filePath,
		lock:     &sync.Mutex{},
	}
}

func (w *worker) initFile(ctx context.Context) error {
	isEmpty, err := w.f.IsFileEmpty(w.filePath)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating headers for nodes file", "error", err)
		return err
	}
	if !isEmpty {
		return nil
	}

	w.lock.Lock()
	defer w.lock.Unlock()

	return w.csv.CreateFileWithHeader(ctx, w.filePath, graph.NodeCsvHeader)
}

func (w *worker) ReadNodes(ctx context.Context) ([]graph.Node, error) {
	if err := w.initFile(ctx); err != nil {
		return nil, err
	}

	w.lock.Lock()
	nodes, err := w.csv.ReadNodesFromFile(ctx, w.filePath)
	w.lock.Unlock()

	if err != nil {
		slog.ErrorContext(ctx, "Error reading nodes from CSV file", "filePath", w.filePath, "error", err)
		return nil, err
	}

	if nodes == nil {
		return EmptyGraphNodes, nil
	}

	return nodes, nil
}

func (w *worker) WriteNodes(ctx context.Context, nodes []graph.Node) error {
	if err := w.initFile(ctx); err != nil {
		return err
	}

	w.lock.Lock()
	err := w.csv.WriteNodesAsCsv(ctx, w.filePath, nodes)
	w.lock.Unlock()

	if err != nil {
		slog.ErrorContext(ctx, "Error writing nodes to CSV file", "filePath", w.filePath, "error", err)
	}
	return err
}
