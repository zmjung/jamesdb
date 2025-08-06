package grapher

import (
	"context"
	"log/slog"

	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/disk"
)

const (
	NodeCsvHeader = "id,type,name,edges,traits\n"
)

var EmptyGraphNodes = []graph.Node{}

type GraphService struct {
	StorageRootPath string
	NodePath        string
}

func NewGraphService(cfg *config.Config) *GraphService {
	nodePath, err := disk.AddFolder(cfg.Database.RootPath, "nodes")
	if err != nil {
		slog.Error("Error creating nodes folder", "error", err)
		return nil
	}

	slog.Debug("Set node path", "nodePath", nodePath)

	return &GraphService{
		StorageRootPath: cfg.Database.RootPath,
		NodePath:        nodePath,
	}
}

func (gw *GraphService) GetAllNodesByType(ctx context.Context, nodeType string) ([]graph.Node, error) {
	// This function retrieves all nodes of a specific type from the storage.

	filePath := disk.GetFilePath(gw.NodePath, nodeType+".csv")
	isFileEmpty, err := disk.IsFileEmpty(filePath)
	if err != nil {
		slog.ErrorContext(ctx, "Error checking if file is empty", "filePath", filePath, "error", err)
		return nil, err
	}

	if isFileEmpty {
		return EmptyGraphNodes, nil
	}

	nodes, err := disk.ReadNodesFromFile(filePath)
	if err != nil {
		slog.ErrorContext(ctx, "Error reading nodes from CSV file", "filePath", filePath, "error", err)
		return nil, err
	}

	return nodes, nil
}

func (gw *GraphService) WriteNode(ctx context.Context, node *graph.Node) error {
	// Converts the node data to a csv format
	// and writes it to a disk file.

	filePath, err := prepareFileAndGetPath(ctx, node.Type, gw.NodePath)
	if err != nil {
		return err
	}

	// save the csv data to a file
	// file name is based on node type
	err = disk.WriteNodesAsCsv(filePath, []graph.Node{*node})
	if err != nil {
		slog.ErrorContext(ctx, "Error writing nodes to CSV file", "filePath", filePath, "error", err)
	}
	return err
}

func prepareFileAndGetPath(ctx context.Context, nodeType string, nodePath string) (string, error) {
	filePath := disk.GetFilePath(nodePath, nodeType+".csv")

	isFileEmpty, err := disk.IsFileEmpty(filePath)
	if err != nil {
		slog.ErrorContext(ctx, "Error checking if file is empty", "filePath", filePath, "error", err)
		return "", err
	}

	if isFileEmpty {
		// if the file is empty, setup header
		slog.DebugContext(ctx, "Creating CSV header because file is empty", "header", NodeCsvHeader)
		err = disk.WriteCsvToFile(filePath, NodeCsvHeader)
		if err != nil {
			return "", nil
		}
	}

	return filePath, nil
}
