package help

const (
	LIKE_TRACK_HELP_TEMPLATE = `🤍🎵 Like track(s) in Spotify by ID.

You must set the args or the flag "id" of track(s), album(s) or artist(s) for like.
If you set both args and flag "id", both will be liked.

If you pass the artist ID, spotlike will like all albums released by the artist.
If you pass the album ID, spotlike will like all tracks included in the album.

Usage:
  spotlike like track [flags]

Flags:
  -f, --force       like track(s) without confirming
  -h, --help        help for track
  -i, --id string   ID of the track(s) or the artist(s) or the album for like
  -v, --verbose     print verbose output
`
)
