//go:build integration

package applescript_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/helmedeiros/itunesScript/internal/adapter/applescript"
)

// TestIntegrationStatus exercises the real osascript path against Music.app.
// It is read-only and skips cleanly when Music is not running. Run with:
//
//	make integration
func TestIntegrationStatus(t *testing.T) {
	p := applescript.New()

	status, err := p.Status(context.Background())
	if errors.Is(err, applescript.ErrNotRunning) {
		t.Skip("Music.app is not running")
	}
	require.NoError(t, err)

	assert.GreaterOrEqual(t, status.Volume.Int(), 0)
	assert.LessOrEqual(t, status.Volume.Int(), 100)
}
