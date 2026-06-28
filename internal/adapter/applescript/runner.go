package applescript

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

// language selects the scripting language osascript interprets.
type language int

const (
	appleScript language = iota
	javaScript
)

// runner executes a script through osascript and returns its stdout. It is the
// single seam between this adapter and the operating system, which keeps the
// rest of the package unit-testable.
type runner interface {
	Run(ctx context.Context, lang language, script string) ([]byte, error)
}

// execRunner runs scripts with the real /usr/bin/osascript.
type execRunner struct{}

func (execRunner) Run(ctx context.Context, lang language, script string) ([]byte, error) {
	args := make([]string, 0, 4)
	if lang == javaScript {
		args = append(args, "-l", "JavaScript")
	}
	args = append(args, "-e", script)

	out, err := exec.CommandContext(ctx, "osascript", args...).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("osascript: %w: %s", err, exitErr.Stderr)
		}
		return nil, fmt.Errorf("osascript: %w", err)
	}
	return out, nil
}
