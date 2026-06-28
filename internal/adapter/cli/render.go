// Package cli is the driving adapter that exposes the application as terminal
// commands. It depends on the Controller port, never on a concrete adapter.
package cli

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/helmedeiros/amp/internal/music"
)

// progressBarWidth is how many cells the status progress bar spans.
const progressBarWidth = 30

// RenderStatus formats a status snapshot for humans, one line per fact.
func RenderStatus(s music.Status) string {
	var b strings.Builder

	fmt.Fprintf(&b, "[%s]", s.State)
	if s.HasTrack() {
		fmt.Fprintf(&b, " %s - %s", s.Track.Artist, s.Track.Name)
		if s.Track.Album != "" {
			fmt.Fprintf(&b, " (%s)", s.Track.Album)
		}
		fmt.Fprintf(&b, "\n%s %s %s",
			FormatClock(s.Elapsed),
			ProgressBar(s.Progress(), progressBarWidth),
			FormatClock(s.Track.Duration),
		)
	}
	fmt.Fprintf(&b, "\nvol %d%%", s.Volume.Int())
	if s.Shuffle {
		b.WriteString("  shuffle")
	}
	if s.Repeat != music.RepeatOff {
		fmt.Fprintf(&b, "  repeat %s", s.Repeat)
	}

	return b.String()
}

// statusJSON is the stable machine-readable shape of a status snapshot.
type statusJSON struct {
	State          string     `json:"state"`
	Volume         int        `json:"volume"`
	Shuffle        bool       `json:"shuffle"`
	Repeat         string     `json:"repeat"`
	ElapsedSeconds float64    `json:"elapsed_seconds"`
	Track          *trackJSON `json:"track,omitempty"`
}

type trackJSON struct {
	Name            string  `json:"name"`
	Artist          string  `json:"artist"`
	Album           string  `json:"album"`
	DurationSeconds float64 `json:"duration_seconds"`
}

// RenderStatusJSON formats a status snapshot as a single JSON object.
func RenderStatusJSON(s music.Status) string {
	payload := statusJSON{
		State:          s.State.String(),
		Volume:         s.Volume.Int(),
		Shuffle:        s.Shuffle,
		Repeat:         s.Repeat.String(),
		ElapsedSeconds: s.Elapsed.Seconds(),
	}
	if s.HasTrack() {
		payload.Track = &trackJSON{
			Name:            s.Track.Name,
			Artist:          s.Track.Artist,
			Album:           s.Track.Album,
			DurationSeconds: s.Track.Duration.Seconds(),
		}
	}

	out, _ := json.Marshal(payload)
	return string(out)
}

// ProgressBar renders a fixed-width bar for a fraction in [0, 1]: filled cells
// for the elapsed portion, empty cells for the rest. Out-of-range fractions are
// clamped, and the result is always exactly width cells wide.
func ProgressBar(fraction float64, width int) string {
	if width <= 0 {
		return ""
	}
	fraction = math.Max(0, math.Min(1, fraction))

	filled := int(math.Round(fraction * float64(width)))
	return strings.Repeat("━", filled) + strings.Repeat("─", width-filled)
}

// FormatClock renders a duration as m:ss (minutes are not zero-padded).
func FormatClock(d time.Duration) string {
	total := int(d.Round(time.Second).Seconds())
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}
