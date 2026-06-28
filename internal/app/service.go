// Package app holds the application use cases. It orchestrates the domain and
// the ports and depends on neither the CLI nor any concrete engine adapter.
package app

import (
	"context"
	"fmt"

	"github.com/helmedeiros/itunesScript/internal/music"
	"github.com/helmedeiros/itunesScript/internal/port"
)

// DefaultUnmuteVolume is the level unmute restores when no prior level is known.
const DefaultUnmuteVolume = 50

// Service is the application's entry point for controlling playback. It is the
// single object the driving adapters (CLI, later the TUI) call into.
type Service struct {
	player port.Player
	volume port.VolumeStore
}

// NewService wires the service to a Player and a VolumeStore.
func NewService(player port.Player, volume port.VolumeStore) *Service {
	return &Service{player: player, volume: volume}
}

var _ port.Controller = (*Service)(nil)

// Status reads the current player snapshot.
func (s *Service) Status(ctx context.Context) (music.Status, error) {
	return s.player.Status(ctx)
}

// Play resumes or starts playback.
func (s *Service) Play(ctx context.Context) error { return s.player.Play(ctx) }

// Pause halts playback.
func (s *Service) Pause(ctx context.Context) error { return s.player.Pause(ctx) }

// Toggle flips between playing and paused.
func (s *Service) Toggle(ctx context.Context) error { return s.player.TogglePlayPause(ctx) }

// Stop halts playback and unloads the current track.
func (s *Service) Stop(ctx context.Context) error { return s.player.Stop(ctx) }

// Next advances to the next track.
func (s *Service) Next(ctx context.Context) error { return s.player.Next(ctx) }

// Previous returns to the previous track.
func (s *Service) Previous(ctx context.Context) error { return s.player.Previous(ctx) }

// SetVolume sets an absolute volume, clamped to the valid range, and returns
// the level that was applied.
func (s *Service) SetVolume(ctx context.Context, level int) (music.Volume, error) {
	v := music.NewVolume(level)
	if err := s.player.SetVolume(ctx, v); err != nil {
		return 0, fmt.Errorf("set volume: %w", err)
	}
	return v, nil
}

// SetShuffle enables or disables shuffle.
func (s *Service) SetShuffle(ctx context.Context, enabled bool) error {
	return s.player.SetShuffle(ctx, enabled)
}

// ToggleShuffle flips shuffle relative to its current state and returns the new
// value. The current state is read first; if that read fails, no change is made.
func (s *Service) ToggleShuffle(ctx context.Context) (bool, error) {
	status, err := s.player.Status(ctx)
	if err != nil {
		return false, fmt.Errorf("read shuffle: %w", err)
	}

	enabled := !status.Shuffle
	if err := s.player.SetShuffle(ctx, enabled); err != nil {
		return false, fmt.Errorf("set shuffle: %w", err)
	}
	return enabled, nil
}

// SetRepeat sets the repeat mode.
func (s *Service) SetRepeat(ctx context.Context, mode music.RepeatMode) error {
	return s.player.SetRepeat(ctx, mode)
}

// AdjustVolume shifts the current volume by delta, clamped to the valid range,
// and returns the new level. The current volume is read first; if that read
// fails, no change is applied.
func (s *Service) AdjustVolume(ctx context.Context, delta int) (music.Volume, error) {
	status, err := s.player.Status(ctx)
	if err != nil {
		return 0, fmt.Errorf("read volume: %w", err)
	}

	v := status.Volume.Adjust(delta)
	if err := s.player.SetVolume(ctx, v); err != nil {
		return 0, fmt.Errorf("set volume: %w", err)
	}
	return v, nil
}

// Mute remembers the current volume and sets it to zero. If the volume is
// already zero it does nothing, so the remembered level survives a double mute.
func (s *Service) Mute(ctx context.Context) error {
	status, err := s.player.Status(ctx)
	if err != nil {
		return fmt.Errorf("read volume: %w", err)
	}
	if status.Volume.IsMuted() {
		return nil
	}

	if err := s.volume.Save(status.Volume.Int()); err != nil {
		return fmt.Errorf("remember volume: %w", err)
	}
	if err := s.player.SetVolume(ctx, music.NewVolume(0)); err != nil {
		return fmt.Errorf("mute: %w", err)
	}
	return nil
}

// Unmute restores the remembered volume, falling back to DefaultUnmuteVolume
// when no prior level is known, and returns the level applied.
func (s *Service) Unmute(ctx context.Context) (music.Volume, error) {
	level, ok, err := s.volume.Load()
	if err != nil {
		return 0, fmt.Errorf("recall volume: %w", err)
	}
	if !ok || level <= 0 {
		level = DefaultUnmuteVolume
	}

	v := music.NewVolume(level)
	if err := s.player.SetVolume(ctx, v); err != nil {
		return 0, fmt.Errorf("unmute: %w", err)
	}
	return v, nil
}
