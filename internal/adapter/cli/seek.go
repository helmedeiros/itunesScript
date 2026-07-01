package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/helmedeiros/amp/internal/music"
)

// parseSeekArg interprets a seek argument into a mode and value:
//
//	"90"    -> absolute 90 seconds
//	"1:30"  -> absolute 90 seconds (mm:ss)
//	"+10"   -> relative +10 seconds
//	"-10"   -> relative -10 seconds
//	"50%"   -> 50 percent of the track
func parseSeekArg(arg string) (music.SeekMode, float64, error) {
	s := strings.TrimSpace(arg)
	if s == "" {
		return 0, 0, fmt.Errorf("empty seek argument")
	}

	if pct, ok := strings.CutSuffix(s, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(pct), 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid seek percentage %q", arg)
		}
		return music.SeekPercent, v, nil
	}

	if s[0] == '+' || s[0] == '-' {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid relative seek %q", arg)
		}
		return music.SeekRelative, v, nil
	}

	if strings.Contains(s, ":") {
		secs, err := parseClock(s)
		if err != nil {
			return 0, 0, err
		}
		return music.SeekAbsolute, secs, nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid seek %q: want seconds, mm:ss, +n, -n or n%%", arg)
	}
	return music.SeekAbsolute, v, nil
}

// parseClock parses "mm:ss" into seconds.
func parseClock(s string) (float64, error) {
	m, sec, ok := strings.Cut(s, ":")
	if !ok {
		return 0, fmt.Errorf("invalid time %q", s)
	}
	mins, err1 := strconv.Atoi(strings.TrimSpace(m))
	secs, err2 := strconv.Atoi(strings.TrimSpace(sec))
	if err1 != nil || err2 != nil || secs < 0 || secs >= 60 {
		return 0, fmt.Errorf("invalid time %q: want mm:ss", s)
	}
	return float64(mins*60 + secs), nil
}
