package music_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/helmedeiros/itunesScript/internal/music"
)

func TestParsePlayerState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  music.PlayerState
	}{
		{name: "playing", input: "playing", want: music.Playing},
		{name: "paused", input: "paused", want: music.Paused},
		{name: "stopped", input: "stopped", want: music.Stopped},
		{name: "trims whitespace", input: "  playing\n", want: music.Playing},
		{name: "case insensitive", input: "Playing", want: music.Playing},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := music.ParsePlayerState(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParsePlayerStateUnknown(t *testing.T) {
	t.Parallel()

	_, err := music.ParsePlayerState("fast-forwarding")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fast-forwarding")
}

func TestPlayerStateString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "playing", music.Playing.String())
	assert.Equal(t, "paused", music.Paused.String())
	assert.Equal(t, "stopped", music.Stopped.String())
}

func TestPlayerStateIsPlaying(t *testing.T) {
	t.Parallel()

	assert.True(t, music.Playing.IsPlaying())
	assert.False(t, music.Paused.IsPlaying())
	assert.False(t, music.Stopped.IsPlaying())
}
