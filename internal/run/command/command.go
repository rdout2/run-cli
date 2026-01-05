package command

import (
	"io"

	"github.com/JulienBreux/run-cli/internal/run/command/version"
	"github.com/JulienBreux/run-cli/internal/run/config"
	"github.com/JulienBreux/run-cli/internal/run/tui/app"
	"github.com/spf13/cobra"
)

// New returns the root of CLI.
func New(in io.Reader, out, err io.Writer) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "run",
		Short: "Run is a CLI to play with Google Cloud Run interactively.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			return app.Run(cfg)
		},
	}

	cmd.AddCommand(version.NewCmdVersion(in, out, err))

	return
}
