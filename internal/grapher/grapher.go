package grapher

import (
	"fmt"

	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/disk"
)

const (
	NodeCsvHeader = "id,type,name,edges,traits\n"
)

type GraphService struct {
	StorageRootPath string
	NodePath        string
}

func NewGraphService(cfg *config.Config) *GraphService {
	nodePath, err := disk.AddFolder(cfg.Database.RootPath, "nodes")
	if err != nil {
		fmt.Printf("Error creating nodes folder: %v\n", err)
		return nil
	}

	fmt.Println("Node path:", nodePath)

	return &GraphService{
		StorageRootPath: cfg.Database.RootPath,
		NodePath:        nodePath,
	}
}

func (gw *GraphService) GetAllNodesByType(nodeType string) ([]graph.Node, error) {
	// This function retrieves all nodes of a specific type from the storage.

	filePath := disk.GetFilePath(gw.NodePath, nodeType+".csv")
	nodes, err := disk.ReadNodesFromFile(filePath)
	if err != nil {
		fmt.Printf("Error reading nodes from CSV file: %v\n", err)
		return nil, err
	}

	return nodes, nil
}

func (gw *GraphService) WriteNode(node *graph.Node) error {
	// Converts the node data to a csv format
	// and writes it to a disk file.

	filePath := disk.GetFilePath(gw.NodePath, node.Type+".csv")

	isFileEmpty, err := disk.IsFileEmpty(filePath)
	if err != nil {
		fmt.Printf("Error checking if file is empty: %v\n", err)
		return err
	}

	if isFileEmpty {
		// if the file is empty, setup header
		disk.WriteCsvToFile(filePath, NodeCsvHeader)
	}

	// save the csv data to a file
	// file name is based on node type
	return disk.WriteNodesAsCsv(filePath, []graph.Node{*node})
}
