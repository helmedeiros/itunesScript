package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/helmedeiros/itunesScript/internal/music"
	"github.com/helmedeiros/itunesScript/internal/port"
)

// NewRootCmd builds the command tree wired to a Controller. Output is written
// to the command's configured streams so it can be captured in tests.
func NewRootCmd(ctrl port.Controller) *cobra.Command {
	root := &cobra.Command{
		Use:           "am",
		Short:         "Control Apple Music from the terminal",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		statusCmd(ctrl),
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

func statusCmd(ctrl port.Controller) *cobra.Command {
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
			if asJSON {
				fmt.Fprintln(cmd.OutOrStdout(), RenderStatusJSON(s))
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), RenderStatus(s))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "output machine-readable JSON")

	return cmd
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
