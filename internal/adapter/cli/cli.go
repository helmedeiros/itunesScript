package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"

	"github.com/helmedeiros/amp/internal/music"
	"github.com/helmedeiros/amp/internal/port"
)

// NewRootCmd builds the command tree wired to a Controller. Output is written
// to the command's configured streams so it can be captured in tests.
func NewRootCmd(ctrl port.Controller) *cobra.Command {
	var noColor bool

	root := &cobra.Command{
		Use:           "amp",
		Short:         "Control Apple Music from the terminal",
		Version:       buildVersion(),
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")

	root.AddCommand(
		statusCmd(ctrl, &noColor),
		nowCmd(ctrl),
		searchCmd(ctrl),
		playlistsCmd(ctrl),
		libraryCmd(ctrl),
		transportCmd(ctrl, "open", "Launch Apple Music", port.Controller.Open),
		transportCmd(ctrl, "play", "Resume or start playback", port.Controller.Play),
		transportCmd(ctrl, "pause", "Pause playback", port.Controller.Pause),
		transportCmd(ctrl, "toggle", "Toggle play/pause", port.Controller.Toggle),
		transportCmd(ctrl, "stop", "Stop playback", port.Controller.Stop),
		transportCmd(ctrl, "next", "Skip to the next track", port.Controller.Next),
		transportCmd(ctrl, "prev", "Skip to the previous track", port.Controller.Previous),
		volCmd(ctrl),
		muteCmd(ctrl),
		unmuteCmd(ctrl),
		shuffleCmd(ctrl),
		repeatCmd(ctrl),
	)

	return root
}

// buildVersion reports the module version this binary was built from, falling
// back to "dev" for local builds without version metadata.
func buildVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "dev"
}

func nowCmd(ctrl port.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "now",
		Short: "Print the current track on one line",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := ctrl.Status(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), RenderNow(s))
			return nil
		},
	}
}

func searchCmd(ctrl port.Controller) *cobra.Command {
	var (
		limit  int
		asJSON bool
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search the library",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tracks, err := ctrl.Search(cmd.Context(), strings.Join(args, " "), limit)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if asJSON {
				fmt.Fprintln(out, RenderTracksJSON(tracks))
			} else {
				fmt.Fprintln(out, RenderTracks(tracks))
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 50, "maximum results (0 for all)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "output machine-readable JSON")

	return cmd
}

func libraryCmd(ctrl port.Controller) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "library",
		Short: "Browse the library",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		namesSubcmd("artists", "List all artists", ctrl.Artists),
		namesSubcmd("albums", "List all albums", ctrl.Albums),
	)
	return cmd
}

// namesSubcmd builds a subcommand that prints a list of names fetched from the
// controller, honoring --json.
func namesSubcmd(use, short string, fetch func(context.Context) ([]string, error)) *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			names, err := fetch(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), RenderNames(names, asJSON))
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "output machine-readable JSON")

	return cmd
}

func playlistsCmd(ctrl port.Controller) *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "playlists",
		Short: "List your playlists",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			playlists, err := ctrl.Playlists(cmd.Context())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if asJSON {
				fmt.Fprintln(out, RenderPlaylistsJSON(playlists))
			} else {
				fmt.Fprintln(out, RenderPlaylists(playlists))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "output machine-readable JSON")

	return cmd
}

func muteCmd(ctrl port.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "mute",
		Short: "Silence playback, remembering the current volume",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := ctrl.Mute(cmd.Context()); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "muted")
			return nil
		},
	}
}

func unmuteCmd(ctrl port.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "unmute",
		Short: "Restore the volume from before muting",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			v, err := ctrl.Unmute(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "vol %d%%\n", v.Int())
			return nil
		},
	}
}

func onOff(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
}

func statusCmd(ctrl port.Controller, noColor *bool) *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show what is playing",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			s, err := ctrl.Status(cmd.Context())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if asJSON {
				fmt.Fprintln(out, RenderStatusJSON(s))
				return nil
			}

			theme := PlainTheme
			if wantsColor(out, *noColor) {
				theme = ColorTheme()
			}
			fmt.Fprintln(out, RenderStatus(s, theme))
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "output machine-readable JSON")

	return cmd
}

// wantsColor reports whether colored output should be emitted: never when
// disabled by --no-color or the NO_COLOR convention, and otherwise only when
// the destination is a terminal.
func wantsColor(w io.Writer, noColor bool) bool {
	if noColor || os.Getenv("NO_COLOR") != "" {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// transportCmd builds a no-argument command that invokes a single Controller
// action.
func transportCmd(ctrl port.Controller, use, short string, action func(port.Controller, context.Context) error) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return action(ctrl, cmd.Context())
		},
	}
}

func shuffleCmd(ctrl port.Controller) *cobra.Command {
	return &cobra.Command{
		Use:       "shuffle [on|off|toggle]",
		Short:     "Turn shuffle on, off, or toggle it (default: toggle)",
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"on", "off", "toggle"},
		RunE: func(cmd *cobra.Command, args []string) error {
			action := "toggle"
			if len(args) == 1 {
				action = args[0]
			}

			var (
				enabled bool
				err     error
			)
			switch action {
			case "on":
				enabled, err = true, ctrl.SetShuffle(cmd.Context(), true)
			case "off":
				enabled, err = false, ctrl.SetShuffle(cmd.Context(), false)
			case "toggle":
				enabled, err = ctrl.ToggleShuffle(cmd.Context())
			default:
				return fmt.Errorf("invalid shuffle argument %q: want on, off or toggle", action)
			}
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "shuffle %s\n", onOff(enabled))
			return nil
		},
	}
}

func repeatCmd(ctrl port.Controller) *cobra.Command {
	return &cobra.Command{
		Use:       "repeat <off|one|all>",
		Short:     "Set the repeat mode",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"off", "one", "all"},
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := music.ParseRepeatMode(args[0])
			if err != nil {
				return err
			}
			if err := ctrl.SetRepeat(cmd.Context(), mode); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "repeat %s\n", mode)
			return nil
		},
	}
}

func volCmd(ctrl port.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "vol <n|+n|-n|up|down>",
		Short: "Get or set the volume",
		Args:  cobra.ExactArgs(1),
		// A relative argument like "-20" would otherwise be parsed as flags.
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "-h" || args[0] == "--help" {
				return cmd.Help()
			}

			relative, value, err := parseVolumeArg(args[0])
			if err != nil {
				return err
			}

			var applied music.Volume
			if relative {
				applied, err = ctrl.AdjustVolume(cmd.Context(), value)
			} else {
				applied, err = ctrl.SetVolume(cmd.Context(), value)
			}
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "vol %d%%\n", applied.Int())
			return nil
		},
	}
}
