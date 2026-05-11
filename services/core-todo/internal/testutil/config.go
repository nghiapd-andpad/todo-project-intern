package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
)

func NewTestConfig(t *testing.T) *config.Config {
	t.Helper()

	cfg, err := config.New()
	require.NoError(t, err)

	return cfg
}
