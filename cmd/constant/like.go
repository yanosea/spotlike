package constant

const (
	LIKE_USE           = "like"
	LIKE_HELP_TEMPLATE = `🤍 Like content in Spotify by ID.

You must use sub command below...

  * 🎤 artist
  * 💿 album
  * 🎵 track

Usage:
  spotlike like [flags]
  spotlike like [command]

Available Commands:
  album       🤍💿 Like album(s) in Spotify by ID.
  artist      🤍🎤 Like artist(s) in Spotify by ID.
  track       🤍🎵 Like track(s) in Spotify by ID.

Flags:
  -h, --help   help for like

Use "spotlike like [command] --help" for more information about a command.
`
	LIKE_MESSAGE_NO_SUB_COMMAND = `Use sub command below...

  * 🎤 artist
  * 💿 album
  * 🎵 track`
	LIKE_ERROR_MESSAGE_ARGS_OR_FLAG_ID_REQUIRED = `The arguments or the flag "id" is required...`
)
