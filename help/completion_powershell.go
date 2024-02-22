package help

const (
	COMPLETION_POWERSHELL_HELP_TEMPLATE = `🔧🎭 Generate the autocompletion script for the powershell shell.

To load completions in your current shell session:

  spotlike completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command to your powershell profile.

Usage:
  spotlike completion powershell [flags]

Flags:
  -h, --help   help for powershell
`
)
