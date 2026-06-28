// Package cli is the driving adapter that exposes the application as terminal
// commands. It depends on the Controller port, never on a concrete adapter.
package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/helmedeiros/itunesScript/internal/music"
)

// RenderStatus formats a status snapshot for humans, one line per fact.
func RenderStatus(s music.Status) string {
	var b strings.Builder

	fmt.Fprintf(&b, "[%s]", s.State)
	if s.HasTrack() {
		fmt.Fprintf(&b, " %s - %s", s.Track.Artist, s.Track.Name)
		if s.Track.Album != "" {
			fmt.Fprintf(&b, " (%s)", s.Track.Album)
		}
		fmt.Fprintf(&b, "\n%s/%s", FormatClock(s.Elapsed), FormatClock(s.Track.Duration))
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

// FormatClock renders a duration as m:ss (minutes are not zero-padded).
func FormatClock(d time.Duration) string {
	total := int(d.Round(time.Second).Seconds())
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}
