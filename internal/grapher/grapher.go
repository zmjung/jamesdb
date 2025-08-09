package grapher

import (
	"context"
	"log/slog"
	"sync"

	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/disk"
)

const (
	NodeCsvHeader = "id,type,name,edges,traits\n"
)

var instance Grapher
var grapherOnce sync.Once

type Grapher interface {
	ReadNodesByType(ctx context.Context, nodeType string) ([]graph.Node, error)
	WriteNode(ctx context.Context, node *graph.Node) error
}

type graphService struct {
	StorageRootPath  string
	NodePath         string
	NodeTypeToWorker map[string]worker
	Lock             *sync.Mutex
}

func GetInstance(cfg *config.Config) Grapher {
	grapherOnce.Do(func() {
		instance = newGrapher(cfg)
	})
	return instance
}

func newGrapher(cfg *config.Config) Grapher {
	nodePath, err := disk.AddFolder(cfg.Database.RootPath, "nodes")
	if err != nil {
		slog.Error("Error creating nodes folder", "error", err)
		return nil
	}

	slog.Debug("Set node path", "nodePath", nodePath)

	return &graphService{
		StorageRootPath:  cfg.Database.RootPath,
		NodePath:         nodePath,
		NodeTypeToWorker: make(map[string]worker),
		Lock:             &sync.Mutex{},
	}
}

func (gs *graphService) getWorker(nodeType string) worker {
	w, exists := gs.NodeTypeToWorker[nodeType]
	if exists {
		return w
	}

	gs.Lock.Lock()
	defer gs.Lock.Unlock()

	w = NewWorker(gs, nodeType)
	gs.NodeTypeToWorker[nodeType] = w
	return w
}

func (gs *graphService) ReadNodesByType(ctx context.Context, nodeType string) ([]graph.Node, error) {
	return gs.getWorker(nodeType).ReadNodes(ctx)
}

func (gs *graphService) WriteNode(ctx context.Context, node *graph.Node) error {
	return gs.getWorker(node.Type).WriteNodes(ctx, []graph.Node{*node})
}
