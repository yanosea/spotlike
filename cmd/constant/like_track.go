package constant

const (
	LIKE_TRACK_USE           = "track"
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
	LIKE_TRACK_FLAG_ID                  = "id"
	LIKE_TRACK_FLAG_ID_SHORTHAND        = "i"
	LIKE_TRACK_FLAG_ID_DESCRIPTION      = "ID of the track(s) or the artist(s) or the album for like"
	LIKE_TRACK_FLAG_FORCE               = "force"
	LIKE_TRACK_FLAG_FORCE_SHORTHAND     = "f"
	LIKE_TRACK_FLAG_FORCE_DEFAULT       = false
	LIKE_TRACK_FLAG_FORCE_DESCRIPTION   = "like track(s) without confirming"
	LIKE_TRACK_FLAG_VERBOSE             = "verbose"
	LIKE_TRACK_FLAG_VERBOSE_SHORTHAND   = "v"
	LIKE_TRACK_FLAG_VERBOSE_DEFAULT     = false
	LIKE_TRACK_FLAG_VERBOSE_DESCRIPTION = "print verbose output"

	LIKE_TRACK_MESSAGE_TEMPLATE_LIKE_TRACK_ALREADY_LIKED              = "✨ %s in [%s] by [%s] already liked!\t:\t[%s]"
	LIKE_TRACK_MESSAGE_TEMPLATE_LIKE_TRACK_REFUSED                    = "❌ Like %s in [%s] by [%s] refused!\t:\t[%s]"
	LIKE_TRACK_MESSAGE_TEMPLATE_LIKE_TRACK_SUCCEEDED                  = "✅ Like %s in [%s] by [%s] succeeded!\t:\t[%s]"
	LIKE_TRACK_ERROR_MESSAGE_TEMPLATE_ID_NOT_ARTIST_ALBUM_TRACK       = "The ID you passed [%s] is not track, album or artist..."
	LIKE_TRACK_PROMPT_LABEL_TEMPLATE_CONFIRM_LIKE_ALL_TRACK_BY_ARTIST = "❔ Do you execute like all tracks by [%s]"
	LIKE_TRACK_PROMPT_LABEL_TEMPLATE_CONFIRM_LIKE_ALL_TRACK_IN_ALBUM  = "❔ Do you execute like all tracks in [%s]"
)
