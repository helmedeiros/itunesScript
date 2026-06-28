package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWantsColor(t *testing.T) {
	// Not parallel: mutates the NO_COLOR environment variable.
	t.Setenv("NO_COLOR", "")

	regular, err := os.Create(filepath.Join(t.TempDir(), "out"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = regular.Close() })

	t.Run("disabled by flag", func(t *testing.T) {
		assert.False(t, wantsColor(regular, true))
	})

	t.Run("disabled by NO_COLOR", func(t *testing.T) {
		t.Setenv("NO_COLOR", "1")
		assert.False(t, wantsColor(regular, false))
	})

	t.Run("non-file writer is never a terminal", func(t *testing.T) {
		assert.False(t, wantsColor(&bytes.Buffer{}, false))
	})

	t.Run("regular file is not a terminal", func(t *testing.T) {
		assert.False(t, wantsColor(regular, false))
	})
}
