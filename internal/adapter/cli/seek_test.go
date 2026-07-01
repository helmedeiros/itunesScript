package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/helmedeiros/amp/internal/music"
)

func TestParseSeekArg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		arg      string
		wantMode music.SeekMode
		wantVal  float64
	}{
		{name: "absolute seconds", arg: "90", wantMode: music.SeekAbsolute, wantVal: 90},
		{name: "mm:ss", arg: "1:30", wantMode: music.SeekAbsolute, wantVal: 90},
		{name: "relative plus", arg: "+10", wantMode: music.SeekRelative, wantVal: 10},
		{name: "relative minus", arg: "-10", wantMode: music.SeekRelative, wantVal: -10},
		{name: "percent", arg: "50%", wantMode: music.SeekPercent, wantVal: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mode, val, err := parseSeekArg(tt.arg)

			require.NoError(t, err)
			assert.Equal(t, tt.wantMode, mode)
			assert.InDelta(t, tt.wantVal, val, 0.001)
		})
	}
}

func TestParseSeekArgErrors(t *testing.T) {
	t.Parallel()

	for _, arg := range []string{"", "abc", "1:90", "10x", "%"} {
		_, _, err := parseSeekArg(arg)
		require.Error(t, err, "arg %q should be rejected", arg)
	}
}
