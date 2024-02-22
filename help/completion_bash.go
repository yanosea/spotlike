package help

const (
	COMPLETION_BASH_HELP_TEMPLATE = `🔧🐚 Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

  source <(spotlike completion bash)

To load completions for every new session, execute once:

  * Linux:

    spotlike completion bash > /etc/bash_completion.d/spotlike

  * macOS:

    spotlike completion bash > $(brew --prefix)/etc/bash_completion.d/spotlike

You will need to start a new shell for this setup to take effect.

Usage:
  spotlike completion bash [flags]

Flags:
  -h, --help   help for bash
`
)
