package cli_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/helmedeiros/amp/internal/adapter/cli"
	"github.com/helmedeiros/amp/internal/music"
)

func playingStatus() music.Status {
	return music.Status{
		State:   music.Playing,
		Volume:  music.NewVolume(60),
		Shuffle: true,
		Repeat:  music.RepeatAll,
		Elapsed: 117 * time.Second,
		Track: music.Track{
			Name:     "Gorgon",
			Artist:   "Utsu-P",
			Album:    "Unknown Album",
			Duration: 255 * time.Second,
		},
	}
}

func TestRenderStatusHuman(t *testing.T) {
	t.Parallel()

	got := cli.RenderStatus(playingStatus())

	assert.Contains(t, got, "[playing]")
	assert.Contains(t, got, "Utsu-P - Gorgon")
	assert.Contains(t, got, "Unknown Album")
	assert.Contains(t, got, "1:57") // elapsed
	assert.Contains(t, got, "4:15") // duration
	assert.Contains(t, got, "━")    // progress bar is rendered
	assert.Contains(t, got, "vol 60%")
}

func TestRenderStatusHumanStopped(t *testing.T) {
	t.Parallel()

	got := cli.RenderStatus(music.Status{State: music.Stopped, Volume: music.NewVolume(50)})

	assert.Contains(t, got, "[stopped]")
	assert.NotContains(t, got, " - ") // no "artist - title" line when no track
}

func TestRenderStatusJSON(t *testing.T) {
	t.Parallel()

	got := cli.RenderStatusJSON(playingStatus())

	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &payload))

	assert.Equal(t, "playing", payload["state"])
	assert.Equal(t, float64(60), payload["volume"])
	assert.Equal(t, true, payload["shuffle"])
	assert.Equal(t, "all", payload["repeat"])
	assert.InDelta(t, 117.0, payload["elapsed_seconds"], 0.001)

	track, ok := payload["track"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Gorgon", track["name"])
	assert.Equal(t, "Utsu-P", track["artist"])
	assert.InDelta(t, 255.0, track["duration_seconds"], 0.001)
}

func TestProgressBar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fraction float64
		width    int
		want     string
	}{
		{name: "empty", fraction: 0, width: 5, want: "─────"},
		{name: "full", fraction: 1, width: 5, want: "━━━━━"},
		{name: "half", fraction: 0.5, width: 10, want: "━━━━━─────"},
		{name: "clamps below zero", fraction: -1, width: 4, want: "────"},
		{name: "clamps above one", fraction: 2, width: 4, want: "━━━━"},
		{name: "zero width is empty", fraction: 0.5, width: 0, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, cli.ProgressBar(tt.fraction, tt.width))
		})
	}
}

func TestProgressBarAlwaysRendersWidthCells(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 20, len([]rune(cli.ProgressBar(0.37, 20))))
}

func TestFormatClock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   time.Duration
		want string
	}{
		{in: 0, want: "0:00"},
		{in: 9 * time.Second, want: "0:09"},
		{in: 117 * time.Second, want: "1:57"},
		{in: 255 * time.Second, want: "4:15"},
		{in: 3725 * time.Second, want: "62:05"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, cli.FormatClock(tt.in))
	}
}
