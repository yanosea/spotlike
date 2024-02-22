package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	help_use   = "help"
	help_short = "🤝 Help provides help for any command in the application."
	help_long  = `🤝 Help provides help for any command in the application.

Simply type spotlike help [path to command] for full details.`
)

func newHelpCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   help_use,
		Short: help_short,
		Long:  help_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
}
