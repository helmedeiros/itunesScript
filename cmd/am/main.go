// Command am controls Apple Music from the terminal.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/helmedeiros/itunesScript/internal/adapter/applescript"
	"github.com/helmedeiros/itunesScript/internal/adapter/cli"
	"github.com/helmedeiros/itunesScript/internal/app"
)

func main() {
	svc := app.NewService(applescript.New())
	root := cli.NewRootCmd(svc)

	if err := root.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "am:", err)
		os.Exit(1)
	}
}
