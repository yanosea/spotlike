package help

const (
	LIKE_ARTIST_HELP_TEMPLATE = `🤍🎤  Like artist(s) in Spotify by ID.

You must set the args or the flag "id" of artist(s) for like.
If you set both args and flag "id", both will be liked.

Usage:
  spotlike like artist [flags]

Flags:
  -f, --force       like artist(s) without confirming
  -h, --help        help for artist
  -i, --id string   ID of the artist(s) for like
  -v, --verbose     print verbose output
`
)
