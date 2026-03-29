package completion

import (
	"bytes"

	c "github.com/spf13/cobra"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// NewCompletionZshCommand returns a new instance of the completion zsh command.
func NewCompletionZshCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("zsh")
	cmd.SetAliases([]string{"zs", "z"})
	cmd.SetUsageTemplate(completionZshUsageTemplate)
	cmd.SetHelpTemplate(completionZshHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runCompletionZsh(cmd, output)
		},
	)

	return cmd
}

// runCompletionZsh generates the autocompletion script for the zsh shell.
func runCompletionZsh(cmd *c.Command, output *string) error {
	buf := new(bytes.Buffer)
	err := cmd.Root().GenZshCompletion(buf)
	*output = buf.String()
	return err
}

const (
	// completionZshHelpTemplate is the help template of the completion zsh command.
	completionZshHelpTemplate = `🔧🧙 Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need to enable it.

You can execute the following once:

  echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

  source <(spotlike completion zsh)

To load completions for every new session, execute once:

  - 🐧 Linux:

    spotlike completion zsh > "${fpath[1]}/_spotlike"

  - 🍎 macOS:

    spotlike completion zsh > $(brew --prefix)/share/zsh/site-functions/_spotlike

You will need to start a new shell for this setup to take effect.

` + completionZshUsageTemplate
	// compleitonUsageTemplate is the usage template of the completion zsh command.
	completionZshUsageTemplate = `Usage:
  spotlike completion zsh [flags]

Flags:
  -h, --help  🤝 help for zsh
`
)
