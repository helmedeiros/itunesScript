// Package store provides file-backed implementations of persistence ports.
package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/helmedeiros/amp/internal/port"
)

// File persists a single volume level as plain text at a given path.
type File struct {
	path string
}

// NewFile returns a File store backed by the given path. Parent directories are
// created on Save.
func NewFile(path string) *File {
	return &File{path: path}
}

var _ port.VolumeStore = (*File)(nil)

// Save writes the level to disk, creating parent directories as needed.
func (f *File) Save(level int) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}
	if err := os.WriteFile(f.path, []byte(strconv.Itoa(level)), 0o600); err != nil {
		return fmt.Errorf("write volume state: %w", err)
	}
	return nil
}

// Load reads the persisted level. A missing file is reported as ok=false with
// no error; an unparseable file is an error.
func (f *File) Load() (level int, ok bool, err error) {
	data, err := os.ReadFile(f.path)
	if os.IsNotExist(err) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("read volume state: %w", err)
	}

	level, err = strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, false, fmt.Errorf("parse volume state: %w", err)
	}
	return level, true, nil
}
