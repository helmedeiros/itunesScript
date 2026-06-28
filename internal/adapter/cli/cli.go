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
	)

	return root
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

func volCmd(ctrl port.Controller) *cobra.Command {
	return &cobra.Command{
		Use:   "vol <n|+n|-n|up|down>",
		Short: "Get or set the volume",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
