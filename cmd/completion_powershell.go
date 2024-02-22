package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	completion_powershell_use   = "powershell"
	completion_powershell_short = "🔧🎭 Generate the autocompletion script for the powershell shell."
	completion_powershell_long  = `🔧🎭 Generate the autocompletion script for the powershell shell.

To load completions in your current shell session:

  spotlike completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command to your powershell profile.`
)

func newCompletionPowerShellCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   completion_powershell_use,
		Short: completion_powershell_short,
		Long:  completion_powershell_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenPowerShellCompletion(globalOption.Out)
		},
	}

	return cmd
}
