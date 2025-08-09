package grapher

import (
	"context"
	"log/slog"
	"sync"

	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/disk"
)

var instance Grapher
var grapherOnce sync.Once

type Grapher interface {
	ReadNodesByType(ctx context.Context, nodeType string) ([]graph.Node, error)
	WriteNode(ctx context.Context, node *graph.Node) error
}

type graphService struct {
	f                disk.FileAccessor
	csv              disk.CsvAccessor
	rootPath         string
	nodePath         string
	nodeTypeToWorker map[string]Worker
	lock             *sync.Mutex
}

func GetInstance(cfg *config.Config, f disk.FileAccessor, csv disk.CsvAccessor) Grapher {
	grapherOnce.Do(func() {
		instance = newGrapher(cfg, f, csv)
	})
	return instance
}

func newGrapher(cfg *config.Config, f disk.FileAccessor, csv disk.CsvAccessor) Grapher {
	nodePath, err := f.AddFolder(cfg.Database.RootPath, "nodes")
	if err != nil {
		slog.Error("Error creating nodes folder", "error", err)
		return nil
	}

	slog.Debug("Set node path", "nodePath", nodePath)

	return &graphService{
		f:                f,
		csv:              csv,
		rootPath:         cfg.Database.RootPath,
		nodePath:         nodePath,
		nodeTypeToWorker: make(map[string]Worker),
		lock:             &sync.Mutex{},
	}
}

func (gs *graphService) getWorker(nodeType string) Worker {
	w, exists := gs.nodeTypeToWorker[nodeType]
	if exists {
		return w
	}

	gs.lock.Lock()
	defer gs.lock.Unlock()

	w = newWorker(gs.f, gs.csv, gs.nodePath, nodeType)
	gs.nodeTypeToWorker[nodeType] = w
	return w
}

func (gs *graphService) ReadNodesByType(ctx context.Context, nodeType string) ([]graph.Node, error) {
	return gs.getWorker(nodeType).ReadNodes(ctx)
}

func (gs *graphService) WriteNode(ctx context.Context, node *graph.Node) error {
	return gs.getWorker(node.Type).WriteNodes(ctx, []graph.Node{*node})
}
