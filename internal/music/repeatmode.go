package music

import (
	"fmt"
	"strings"
)

// RepeatMode is the repeat setting Music.app applies to playback.
type RepeatMode int

const (
	// RepeatOff plays through the queue once.
	RepeatOff RepeatMode = iota
	// RepeatOne repeats the current track.
	RepeatOne
	// RepeatAll repeats the whole queue.
	RepeatAll
)

// String returns the lowercase name, matching the value Music.app uses.
func (m RepeatMode) String() string {
	switch m {
	case RepeatOff:
		return "off"
	case RepeatOne:
		return "one"
	case RepeatAll:
		return "all"
	default:
		return "unknown"
	}
}

// ParseRepeatMode converts a repeat-mode string, as emitted by Music.app, into
// a RepeatMode. Surrounding whitespace is ignored and matching is
// case-insensitive.
func ParseRepeatMode(raw string) (RepeatMode, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "off":
		return RepeatOff, nil
	case "one":
		return RepeatOne, nil
	case "all":
		return RepeatAll, nil
	default:
		return RepeatOff, fmt.Errorf("unknown repeat mode %q", raw)
	}
}
