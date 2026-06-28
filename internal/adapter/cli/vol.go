package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// parseVolumeArg interprets a volume argument. It reports whether the change is
// relative (a delta) or absolute (a target level), along with the value. The
// keywords "up" and "down" map to ±10; a leading "+" or "-" is a delta; a bare
// number is an absolute level.
func parseVolumeArg(arg string) (relative bool, value int, err error) {
	s := strings.TrimSpace(arg)
	switch s {
	case "":
		return false, 0, fmt.Errorf("empty volume argument")
	case "up":
		return true, 10, nil
	case "down":
		return true, -10, nil
	}

	rel := s[0] == '+' || s[0] == '-'
	n, err := strconv.Atoi(s)
	if err != nil {
		return false, 0, fmt.Errorf("invalid volume %q: want a number, +N, -N, up or down", arg)
	}
	return rel, n, nil
}
