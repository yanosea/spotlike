package completion

import (
	"bytes"

	c "github.com/spf13/cobra"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// NewCompletionPowerShellCommand returns a new instance of the completion powershell command.
func NewCompletionPowerShellCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("powershell")
	cmd.SetAliases([]string{"pwsh", "p"})
	cmd.SetUsageTemplate(completionPowerShellUsageTemplate)
	cmd.SetHelpTemplate(completionPowerShellHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runCompletionPowerShell(cmd, output)
		},
	)

	return cmd
}

// runCompletionPowerShell generates the autocompletion script for the powershell shell.
func runCompletionPowerShell(cmd *c.Command, output *string) error {
	buf := new(bytes.Buffer)
	err := cmd.Root().GenPowerShellCompletion(buf)
	*output = buf.String()
	return err
}

const (
	// completionPowerShellHelpTemplate is the help template of the completion powershell command.
	completionPowerShellHelpTemplate = `🔧🪟 Generate the autocompletion script for the powershell shell.

To load completions in your current shell session:

  spotlike completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command to your powershell profile.

` + completionPowerShellUsageTemplate
	// compleitonUsageTemplate is the usage template of the completion powershell command.
	completionPowerShellUsageTemplate = `Usage:
  spotlike completion powershell [flags]

Flags:
  -h, --help  🤝 help for powershell
`
)
