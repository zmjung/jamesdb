package disk

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/jszwec/csvutil"
	"github.com/zmjung/jamesdb/graph"
)

type CsvAccessor interface {
	ReadNodesFromFile(cxt context.Context, filePath string) ([]graph.Node, error)
	WriteNodesAsCsv(cxt context.Context, filePath string, nodes []graph.Node) error
	CreateFileWithHeader(ctx context.Context, filePath string, csvHeader string) error
}

type csvService struct {
	f FileAccessor
}

func NewCsvAccessor(f FileAccessor) CsvAccessor {
	return &csvService{
		f: f,
	}
}

func (s *csvService) ReadNodesFromFile(ctx context.Context, filePath string) ([]graph.Node, error) {
	slog.InfoContext(ctx, "Reading from file", "filePath", filePath)

	reader, err := s.f.GetFileReader(filePath)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	return s.readNodesFromReader(ctx, reader)
}

func (s *csvService) readNodesFromReader(ctx context.Context, reader io.Reader) ([]graph.Node, error) {
	csvReader := csv.NewReader(reader)

	// Preallocate nodes slice for efficiency if possible
	var nodes []graph.Node
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		return nil, err
	}
	if dec == nil {
		err := errors.New("failed to create CSV decoder")
		slog.ErrorContext(ctx, "Decoder cannot be nil", "error", err)
		return nil, err
	}
	dec.Tag = "json" // Use JSON tags for decoding
	dec.WithUnmarshalers(
		csvutil.NewUnmarshalers(
			csvutil.UnmarshalFunc(decodeList),
			csvutil.UnmarshalFunc(decodeMap),
		),
	)
	if err := dec.Decode(&nodes); err != nil && err != io.EOF {
		return nil, err
	}
	return nodes, nil
}

func (s *csvService) WriteNodesAsCsv(ctx context.Context, filePath string, nodes []graph.Node) error {
	slog.InfoContext(ctx, "Writing to file", "filePath", filePath)

	// Create a CSV writer
	writer, err := s.f.GetFileWriter(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()

	return s.writeCsvToWriter(ctx, writer, nodes)
}

func (s *csvService) writeCsvToWriter(ctx context.Context, writer io.Writer, nodes []graph.Node) error {
	// This function writes a slice of graph.Node to a CSV writer.
	csvWriter := csv.NewWriter(writer)

	// Create an encoder
	encoder := csvutil.NewEncoder(csvWriter)
	encoder.AutoHeader = false

	encoder.WithMarshalers(
		csvutil.NewMarshalers(
			csvutil.MarshalFunc(encodeList),
			csvutil.MarshalFunc(encodeMap),
		),
	)

	for _, node := range nodes {
		if err := encoder.Encode(node); err != nil {
			slog.ErrorContext(ctx, "Failed to encode node", "error", err)
			return err
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		slog.ErrorContext(ctx, "Failed to write CSV", "error", err)
		return err
	}
	return nil
}

func (s *csvService) WriteCsvToFile(ctx context.Context, filePath string, csvLine string) error {
	slog.InfoContext(ctx, "Writing to file", "filePath", filePath, "csv", csvLine)

	writer, err := s.f.GetFileWriter(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = writer.Write([]byte(csvLine))
	return err
}

func (s *csvService) CreateFileWithHeader(ctx context.Context, filePath string, csvHeader string) error {
	isFileEmpty, err := s.f.IsFileEmpty(filePath)
	if err != nil {
		return err
	}

	if isFileEmpty {
		// if the file is empty, setup header
		err = s.WriteCsvToFile(ctx, filePath, csvHeader)
	}

	return err
}
