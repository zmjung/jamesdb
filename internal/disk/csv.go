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

func ReadCsv(ctx context.Context, r io.Reader, v any) error {
	csvReader := csv.NewReader(r)

	// Preallocate nodes slice for efficiency if possible
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		return err
	}
	if dec == nil {
		err := errors.New("failed to create CSV decoder")
		slog.ErrorContext(ctx, "Decoder cannot be nil", "error", err)
		return err
	}
	dec.Tag = "json" // Use JSON tags for decoding
	dec.WithUnmarshalers(
		csvutil.NewUnmarshalers(
			csvutil.UnmarshalFunc(decodeList),
			csvutil.UnmarshalFunc(decodeMap),
		),
	)
	if err := dec.Decode(v); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func WriteCsv(ctx context.Context, w io.Writer, v any) error {
	csvWriter := csv.NewWriter(w)

	// Create an encoder
	encoder := csvutil.NewEncoder(csvWriter)
	encoder.AutoHeader = false

	encoder.WithMarshalers(
		csvutil.NewMarshalers(
			csvutil.MarshalFunc(encodeList),
			csvutil.MarshalFunc(encodeMap),
		),
	)

	if err := encoder.Encode(v); err != nil {
		slog.ErrorContext(ctx, "Failed to encode data", "error", err)
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		slog.ErrorContext(ctx, "Failed to write CSV", "error", err)
		return err
	}
	return nil
}

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

	var nodes []graph.Node
	err = ReadCsv(ctx, reader, &nodes)
	return nodes, err
}

func (s *csvService) WriteNodesAsCsv(ctx context.Context, filePath string, nodes []graph.Node) error {
	slog.InfoContext(ctx, "Writing to file", "filePath", filePath)

	// Create a CSV writer
	writer, err := s.f.GetFileWriter(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()

	return WriteCsv(ctx, writer, nodes)
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
