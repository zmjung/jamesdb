package grapher

import (
	"context"
	"log/slog"
	"sync"

	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/disk"
)

var EmptyGraphNodes = []graph.Node{}

type worker interface {
	ReadNodes(ctx context.Context) ([]graph.Node, error)
	WriteNodes(ctx context.Context, nodes []graph.Node) error
}

type impl struct {
	parent   *graphService
	nodeType string
	filePath string
	lock     *sync.Mutex
}

func NewWorker(parent *graphService, nodeType string) worker {
	filePath := disk.GetFilePath(parent.NodePath, nodeType+".csv")
	err := disk.InitFileWithHeader(filePath, graph.NodeCsvHeader)
	if err != nil {
		slog.Error("Error creating headers for nodes file", "error", err)
		return nil
	}

	return &impl{
		parent:   parent,
		nodeType: nodeType,
		filePath: filePath,
		lock:     &sync.Mutex{},
	}
}

func (w *impl) initFile(ctx context.Context) error {
	isEmpty, err := disk.IsFileEmpty(w.filePath)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating headers for nodes file", "error", err)
		return err
	}
	if !isEmpty {
		return nil
	}

	w.lock.Lock()
	defer w.lock.Unlock()

	return disk.InitFileWithHeader(w.filePath, graph.NodeCsvHeader)
}

func (w *impl) ReadNodes(ctx context.Context) ([]graph.Node, error) {
	if err := w.initFile(ctx); err != nil {
		return nil, err
	}

	w.lock.Lock()
	nodes, err := disk.ReadNodesFromFile(w.filePath)
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

func (w *impl) WriteNodes(ctx context.Context, nodes []graph.Node) error {
	if err := w.initFile(ctx); err != nil {
		return err
	}

	w.lock.Lock()
	err := disk.WriteNodesAsCsv(w.filePath, nodes)
	w.lock.Unlock()

	if err != nil {
		slog.ErrorContext(ctx, "Error writing nodes to CSV file", "filePath", w.filePath, "error", err)
	}
	return err
}
