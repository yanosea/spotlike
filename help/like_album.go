package help

const (
	LIKE_ALBUM_HELP_TEMPLATE = `🤍💿 Like album(s) in Spotify by ID.

You must set the args or the flag "id" of album(s) or artist(s) for like.
If you set both args and flag "id", both will be liked.

If you pass the artist ID, spotlike will like all albums released by the artist.

Usage:
  spotlike like album [flags]

Flags:
  -f, --force       like album(s) without confirming
  -h, --help        help for album
  -i, --id string   ID of the album(s) or the artist(s) for like
  -v, --verbose     print verbose output
`
)
