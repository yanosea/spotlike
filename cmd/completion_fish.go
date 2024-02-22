package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	completion_fish_use   = "fish"
	completion_fish_short = "🔧🐟 Generate the autocompletion script for the fish shell."
	completion_fish_long  = `🔧🐟 Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

  spotlike completion fish | source

To load completions for every new session, execute once:

  spotlike completion fish > ~/.config/fish/completions/spotlike.fish

You will need to start a new shell for this setup to take effect.`
)

func newCompletionFishCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   completion_fish_use,
		Short: completion_fish_short,
		Long:  completion_fish_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenFishCompletion(globalOption.Out, false)
		},
	}

	return cmd
}
