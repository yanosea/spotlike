package constant

const (
	ROOT_USE           = "spotlike"
	ROOT_HELP_TEMPLATE = `⚪ spotlike is the CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.

Usage:
  spotlike [flags]
  spotlike [command]

Available Commands:
  auth        🔑 Authenticate your Spotify client.
  completion  🔧 Generate the autocompletion script for the specified shell.
  help        🤝 Help about any command.
  like        🤍 Like content in Spotify by ID.
  search      🔍 Search for the ID of content in Spotify.
  version     🔖 Show the version of spotlike.

Flags:
  -h, --help      🤝 help for spotlike
  -v, --version   🔖 version for spotlike

Use "spotlike [command] --help" for more information about a command.
`
	ROOT_MESSAGE_NO_SUB_COMMAND = `Use sub command below...

  - 🔑 auth
  - 🤍 like
    - 🎤  artist
    - 💿 album
    - 🎵 track
  - 🔍 search
  - 🔖 version

Use "spotlike --help" for more information about spotlike.`
)
