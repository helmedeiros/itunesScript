package music_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/helmedeiros/amp/internal/music"
)

func TestParseRepeatMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  music.RepeatMode
	}{
		{name: "off", input: "off", want: music.RepeatOff},
		{name: "one", input: "one", want: music.RepeatOne},
		{name: "all", input: "all", want: music.RepeatAll},
		{name: "trims and lowercases", input: " All\n", want: music.RepeatAll},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := music.ParseRepeatMode(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseRepeatModeUnknown(t *testing.T) {
	t.Parallel()

	_, err := music.ParseRepeatMode("shuffle")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "shuffle")
}

func TestRepeatModeString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "off", music.RepeatOff.String())
	assert.Equal(t, "one", music.RepeatOne.String())
	assert.Equal(t, "all", music.RepeatAll.String())
}
