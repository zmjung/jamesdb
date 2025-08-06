package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewRequestContext(t *testing.T) {
	ctx := NewRequestContext(context.Background(), map[string]string{
		"key": "value",
	})

	request := ctx.Value(requestKey).(map[string]string)
	require.NotEmpty(t, request)
	require.Equal(t, "value", request["key"])

	newCtx := NewRequestContext(ctx, map[string]string{
		"key": "anotherValue",
	})
	newRequest := newCtx.Value(requestKey).(map[string]string)
	require.NotEmpty(t, newRequest)
	require.Equal(t, "anotherValue", newRequest["key"])
}
