package music_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/helmedeiros/amp/internal/music"
)

func TestNewVolumeClampsToRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "in range", in: 42, want: 42},
		{name: "min", in: 0, want: 0},
		{name: "max", in: 100, want: 100},
		{name: "below min clamps to 0", in: -10, want: 0},
		{name: "above max clamps to 100", in: 250, want: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, music.NewVolume(tt.in).Int())
		})
	}
}

func TestVolumeAdjustClamps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		start int
		delta int
		want  int
	}{
		{name: "increase", start: 50, delta: 10, want: 60},
		{name: "decrease", start: 50, delta: -10, want: 40},
		{name: "increase past max", start: 95, delta: 10, want: 100},
		{name: "decrease past min", start: 5, delta: -10, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := music.NewVolume(tt.start).Adjust(tt.delta)

			assert.Equal(t, tt.want, got.Int())
		})
	}
}

func TestVolumeIsMuted(t *testing.T) {
	t.Parallel()

	assert.True(t, music.NewVolume(0).IsMuted())
	assert.False(t, music.NewVolume(1).IsMuted())
}
