package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_loadConfig(t *testing.T) {
	cfg := loadConfig()
	require.NotNil(t, cfg)
}
