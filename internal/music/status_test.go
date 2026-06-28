package music_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/helmedeiros/amp/internal/music"
)

func TestTrackIsZero(t *testing.T) {
	t.Parallel()

	assert.True(t, music.Track{}.IsZero())
	assert.False(t, music.Track{Name: "Gorgon"}.IsZero())
}

func TestStatusHasTrack(t *testing.T) {
	t.Parallel()

	playing := music.Status{
		State: music.Playing,
		Track: music.Track{Name: "Gorgon", Artist: "Utsu-P"},
	}
	assert.True(t, playing.HasTrack())

	stopped := music.Status{State: music.Stopped}
	assert.False(t, stopped.HasTrack())
}

func TestStatusProgress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		elapsed  time.Duration
		duration time.Duration
		want     float64
	}{
		{name: "half", elapsed: 2 * time.Minute, duration: 4 * time.Minute, want: 0.5},
		{name: "start", elapsed: 0, duration: 4 * time.Minute, want: 0},
		{name: "end", elapsed: 4 * time.Minute, duration: 4 * time.Minute, want: 1},
		{name: "no duration is zero", elapsed: time.Minute, duration: 0, want: 0},
		{name: "overrun clamps to one", elapsed: 5 * time.Minute, duration: 4 * time.Minute, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := music.Status{
				Elapsed: tt.elapsed,
				Track:   music.Track{Duration: tt.duration},
			}

			assert.InDelta(t, tt.want, s.Progress(), 0.0001)
		})
	}
}
