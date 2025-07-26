package grapher

import (
	"fmt"

	"github.com/zmjung/jamesdb/config"
	"github.com/zmjung/jamesdb/graph"
	"github.com/zmjung/jamesdb/internal/converter"
	"github.com/zmjung/jamesdb/internal/disk"
)

const (
	NodeCsvHeader = "ID,Type,Name,Edges,Traits\n"
)

type GraphWriter struct {
	StorageRootPath string
	NodePath        string
}

func NewGraphWriter(cfg *config.Config) *GraphWriter {
	nodePath, err := disk.AddFolder(cfg.Database.RootPath, "nodes")
	if err != nil {
		fmt.Printf("Error creating nodes folder: %v\n", err)
		return nil
	}

	fmt.Println("Node path:", nodePath)

	return &GraphWriter{
		StorageRootPath: cfg.Database.RootPath,
		NodePath:        nodePath,
	}
}

func (gw *GraphWriter) WriteNodeData(node *graph.Node) error {
	// Converts the node data to a csv format
	// and writes it to a disk file.

	// convert node to csv format
	csvData := converter.ConvertToCSV(node)

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
	return disk.WriteCsvToFile(filePath, csvData)
}
