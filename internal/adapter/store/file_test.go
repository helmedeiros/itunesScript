package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/helmedeiros/amp/internal/adapter/store"
	"github.com/helmedeiros/amp/internal/port"
)

func TestFileStoreRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "nested", "volume")
	s := store.NewFile(path)

	require.NoError(t, s.Save(73))

	level, ok, err := s.Load()
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 73, level)
}

func TestFileStoreLoadMissingIsNotAnError(t *testing.T) {
	t.Parallel()

	s := store.NewFile(filepath.Join(t.TempDir(), "absent"))

	_, ok, err := s.Load()

	require.NoError(t, err)
	assert.False(t, ok)
}

func TestFileStoreLoadCorruptReturnsError(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "volume")
	require.NoError(t, os.WriteFile(path, []byte("not-a-number"), 0o600))

	s := store.NewFile(path)

	_, _, err := s.Load()
	require.Error(t, err)
}

func TestFileStoreSatisfiesPort(t *testing.T) {
	t.Parallel()

	var _ port.VolumeStore = store.NewFile("x")
}
