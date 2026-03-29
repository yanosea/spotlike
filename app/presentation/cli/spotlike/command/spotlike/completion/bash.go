package completion

import (
	"bytes"

	c "github.com/spf13/cobra"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// NewCompletionBashCommand returns a new instance of the completion bash command.
func NewCompletionBashCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("bash")
	cmd.SetAliases([]string{"ba", "b"})
	cmd.SetUsageTemplate(completionBashUsageTemplate)
	cmd.SetHelpTemplate(completionBashHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runCompletionBash(cmd, output)
		},
	)

	return cmd
}

// runCompletionBash generates the autocompletion script for the bash shell.
func runCompletionBash(cmd *c.Command, output *string) error {
	buf := new(bytes.Buffer)
	err := cmd.Root().GenBashCompletion(buf)
	*output = buf.String()
	return err
}

const (
	// completionBashHelpTemplate is the help template of the completion bash command.
	completionBashHelpTemplate = `🔧🐚 Generate the autocompletion script for the bash shell.

This script depends on the "bash-completion" package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

  source <(spotlike completion bash)

To load completions for every new session, execute once:

  - 🐧 Linux:

    spotlike completion bash > /etc/bash_completion.d/spotlike

  - 🍎 macOS:

    spotlike completion bash > $(brew --prefix)/etc/bash_completion.d/spotlike

You will need to start a new shell for this setup to take effect.

` + completionBashUsageTemplate
	// compleitonUsageTemplate is the usage template of the completion bash command.
	completionBashUsageTemplate = `Usage:
  spotlike completion bash [flags]

Flags:
  -h, --help  🤝 help for bash
`
)
