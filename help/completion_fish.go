package help

const (
	COMPLETION_FISH_HELP_TEMPLATE = `🔧🐟 Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

  spotlike completion fish | source

To load completions for every new session, execute once:

  spotlike completion fish > ~/.config/fish/completions/spotlike.fish

You will need to start a new shell for this setup to take effect.

Usage:
  spotlike completion fish [flags]

Flags:
  -h, --help   help for fish
`
)
