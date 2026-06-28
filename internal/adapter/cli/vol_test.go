package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVolumeArg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		arg          string
		wantRelative bool
		wantValue    int
	}{
		{name: "absolute", arg: "42", wantRelative: false, wantValue: 42},
		{name: "explicit plus", arg: "+10", wantRelative: true, wantValue: 10},
		{name: "explicit minus", arg: "-10", wantRelative: true, wantValue: -10},
		{name: "up keyword", arg: "up", wantRelative: true, wantValue: 10},
		{name: "down keyword", arg: "down", wantRelative: true, wantValue: -10},
		{name: "whitespace tolerated", arg: " 30 ", wantRelative: false, wantValue: 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rel, val, err := parseVolumeArg(tt.arg)

			require.NoError(t, err)
			assert.Equal(t, tt.wantRelative, rel)
			assert.Equal(t, tt.wantValue, val)
		})
	}
}

func TestParseVolumeArgErrors(t *testing.T) {
	t.Parallel()

	for _, arg := range []string{"", "loud", "+", "12x"} {
		_, _, err := parseVolumeArg(arg)
		require.Error(t, err, "arg %q should be rejected", arg)
	}
}
