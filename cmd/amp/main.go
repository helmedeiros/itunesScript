// Command am controls Apple Music from the terminal.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/helmedeiros/amp/internal/adapter/applescript"
	"github.com/helmedeiros/amp/internal/adapter/cli"
	"github.com/helmedeiros/amp/internal/adapter/store"
	"github.com/helmedeiros/amp/internal/app"
)

func main() {
	svc := app.NewService(applescript.New(), store.NewFile(volumeStatePath()))
	root := cli.NewRootCmd(svc)

	if err := root.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "amp:", err)
		os.Exit(1)
	}
}

// volumeStatePath returns where the pre-mute volume is remembered, under the
// user config dir, falling back to the working directory if it is unavailable.
func volumeStatePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ".am-volume"
	}
	return filepath.Join(dir, "amp", "volume")
}
