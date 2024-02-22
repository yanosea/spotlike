package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	completion_zsh_use   = "zsh"
	completion_zsh_short = "🔧🧙 Generate the autocompletion script for the zsh shell."
	completion_zsh_long  = `🔧🧙 Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need to enable it.

You can execute the following once:

  echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

  source <(spotlike completion zsh)

To load completions for every new session, execute once:

  * Linux:

    spotlike completion zsh > "${fpath[1]}/_spotlike"

  * macOS:

    spotlike completion zsh > $(brew --prefix)/share/zsh/site-functions/_spotlike

You will need to start a new shell for this setup to take effect.`
)

func newCompletionZshCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   completion_zsh_use,
		Short: completion_zsh_short,
		Long:  completion_zsh_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenZshCompletion(globalOption.Out)
		},
	}

	return cmd
}
