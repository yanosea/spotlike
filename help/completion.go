package help

const (
	COMPLETION_HELP_TEMPLATE = `🔧 Generate the autocompletion script for the specified shell.

Usage:
  spotlike completion [flags]
  spotlike completion [command]

Available Commands:
  bash        🔧🐚 Generate the autocompletion script for the bash shell.
  fish        🔧🐟 Generate the autocompletion script for the fish shell.
  powershell  🔧🪟 Generate the autocompletion script for the powershell shell.
  zsh         🔧🧙 Generate the autocompletion script for the zsh shell.

Flags:
  -h, --help   help for completion

Use "spotlike completion [command] --help" for more information about a command.
`
)
