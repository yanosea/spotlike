package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	completion_bash_use   = "bash"
	completion_bash_short = "🔧🐚 Generate the autocompletion script for the bash shell."
	completion_bash_long  = `🔧🐚 Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

  source <(spotlike completion bash)

To load completions for every new session, execute once:

* Linux:

  spotlike completion bash > /etc/bash_completion.d/spotlike

* macOS:

  spotlike completion bash > $(brew --prefix)/etc/bash_completion.d/spotlike

You will need to start a new shell for this setup to take effect.`
)

func newCompletionBashCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   completion_bash_use,
		Short: completion_bash_short,
		Long:  completion_bash_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenBashCompletion(globalOption.Out)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
}
