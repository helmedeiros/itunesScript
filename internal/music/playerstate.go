// Package music holds the pure domain model for controlling Apple Music:
// entities and value objects with no I/O and no external dependencies.
package music

import (
	"fmt"
	"strings"
)

// PlayerState is the transport state Music.app reports for the player.
type PlayerState int

const (
	// Stopped means no track is loaded for playback.
	Stopped PlayerState = iota
	// Playing means a track is currently advancing.
	Playing
	// Paused means a track is loaded but halted.
	Paused
)

// String returns the lowercase name, matching the value Music.app uses.
func (s PlayerState) String() string {
	switch s {
	case Playing:
		return "playing"
	case Paused:
		return "paused"
	case Stopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// IsPlaying reports whether the player is actively advancing a track.
func (s PlayerState) IsPlaying() bool {
	return s == Playing
}

// ParsePlayerState converts a player-state string, as emitted by Music.app,
// into a PlayerState. Surrounding whitespace is ignored and matching is
// case-insensitive.
func ParsePlayerState(raw string) (PlayerState, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "playing":
		return Playing, nil
	case "paused":
		return Paused, nil
	case "stopped":
		return Stopped, nil
	default:
		return Stopped, fmt.Errorf("unknown player state %q", raw)
	}
}
