package graph

import (
	"strings"
	"testing"

	"github.com/jszwec/csvutil"
	"github.com/stretchr/testify/require"
)

func TestNodeType(t *testing.T) {
	header, err := csvutil.Header(Node{}, "json")
	require.NoError(t, err)

	headerCsv := strings.Join(header, ",") + "\n"

	require.Equal(t, headerCsv, NodeCsvHeader)
}
