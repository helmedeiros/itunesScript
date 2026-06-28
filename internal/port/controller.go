package port

import (
	"context"

	"github.com/helmedeiros/amp/internal/music"
)

// Controller is the driving port: the use-case surface that driving adapters
// (the CLI, later the TUI) depend on. The application's Service implements it.
type Controller interface {
	Status(ctx context.Context) (music.Status, error)

	Play(ctx context.Context) error
	Pause(ctx context.Context) error
	Toggle(ctx context.Context) error
	Stop(ctx context.Context) error
	Next(ctx context.Context) error
	Previous(ctx context.Context) error

	// SetVolume sets an absolute level and returns the applied volume.
	SetVolume(ctx context.Context, level int) (music.Volume, error)
	// AdjustVolume shifts the level by delta and returns the new volume.
	AdjustVolume(ctx context.Context, delta int) (music.Volume, error)

	// SetShuffle enables or disables shuffle.
	SetShuffle(ctx context.Context, enabled bool) error
	// ToggleShuffle flips shuffle and returns the new value.
	ToggleShuffle(ctx context.Context) (bool, error)
	// SetRepeat sets the repeat mode.
	SetRepeat(ctx context.Context, mode music.RepeatMode) error

	// Mute silences playback, remembering the current level.
	Mute(ctx context.Context) error
	// Unmute restores the remembered level and returns the applied volume.
	Unmute(ctx context.Context) (music.Volume, error)
}
