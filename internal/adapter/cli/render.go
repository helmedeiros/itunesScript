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

// labelWidth aligns the value column; the longest label is "shuffle" (7).
const labelWidth = 7

type statusRow struct{ label, value string }

// RenderStatus formats a status snapshot for humans as aligned label/value
// rows. Playing and paused show the same fields (only the state word differs);
// stopped/no-track collapses to the state and volume.
func RenderStatus(s music.Status) string {
	var rows []statusRow

	if s.HasTrack() {
		rows = append(rows, statusRow{s.State.String(), artistTitle(s.Track)})
		if s.Track.Album != "" {
			rows = append(rows, statusRow{"album", s.Track.Album})
		}
		rows = append(rows,
			statusRow{"time", timeLine(s)},
			statusRow{"volume", fmt.Sprintf("%d%%", s.Volume.Int())},
			statusRow{"shuffle", onOff(s.Shuffle)},
			statusRow{"repeat", s.Repeat.String()},
		)
	} else {
		rows = append(rows,
			statusRow{s.State.String(), ""},
			statusRow{"volume", fmt.Sprintf("%d%%", s.Volume.Int())},
		)
	}

	lines := make([]string, len(rows))
	for i, r := range rows {
		lines[i] = strings.TrimRight(fmt.Sprintf("%-*s  %s", labelWidth, r.label, r.value), " ")
	}
	return strings.Join(lines, "\n")
}

// artistTitle joins artist and title, or returns the title alone when the
// artist is unknown.
func artistTitle(t music.Track) string {
	if t.Artist == "" {
		return t.Name
	}
	return t.Artist + " — " + t.Name
}

// timeLine renders "elapsed / duration  <bar>  NN%". When the duration is
// unknown it shows a placeholder and omits the bar and percentage.
func timeLine(s music.Status) string {
	if s.Track.Duration <= 0 {
		return FormatClock(s.Elapsed) + " / --:--"
	}
	return fmt.Sprintf("%s / %s  %s  %d%%",
		FormatClock(s.Elapsed),
		FormatClock(s.Track.Duration),
		ProgressBar(s.Progress(), progressBarWidth),
		int(math.Round(s.Progress()*100)),
	)
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
	return fmt.Sprintf("%02d:%02d", total/60, total%60)
}
