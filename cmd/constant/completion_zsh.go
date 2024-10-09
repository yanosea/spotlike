package constant

const (
	COMPLETION_ZSH_USE           = "zsh"
	COMPLETION_ZSH_HELP_TEMPLATE = `🔧🧙 Generate the autocompletion script for the zsh shell.

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

Usage:
  spotlike completion zsh [flags]

Flags:
  -h, --help  🤝 help for zsh
`
)
