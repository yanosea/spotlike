package constant

const (
	LIKE_ALBUM_USE           = "album"
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
	LIKE_ALBUM_FLAG_ID                                    = "id"
	LIKE_ALBUM_SHORTHAND_ID                               = "i"
	LIKE_ALBUM_FLAG_DESCRIPTION_ID                        = "ID of the album(s) or the artist(s) for like"
	LIKE_ALBUM_FLAG_FORCE                                 = "force"
	LIKE_ALBUM_SHORTHAND_FORCE                            = "f"
	LIKE_ALBUM_FLAG_DESCRIPTION_FORCE                     = "like album(s) without confirming"
	LIKE_ALBUM_FLAG_VERBOSE                               = "verbose"
	LIKE_ALBUM_SHORTHAND_VERBOSE                          = "v"
	LIKE_ALBUM_FLAG_DESCRIPTION_VERBOSE                   = "print verbose output"
	LIKE_ALBUM_ERROR_MESSAGE_TEMPLATE_ID_NOT_ARTIST_ALBUM = "The ID you passed [%s] is not album or artist..."
	LIKE_ALBUM_INPUT_LABEL_TEMPLATE_ALL_ALBUM_BY_ARTIST   = "❔ Do you execute like all albums by [%s]"
	LIKE_ALBUM_MESSAGE_TEMPLATE_LIKE_ALBUM_ALREADY_LIKED  = "✨ %s by [%s] already liked!\t:\t[%s]"
	LIKE_ALBUM_MESSAGE_TEMPLATE_LIKE_ALBUM_REFUSED        = "❌ Like %s by [%s] refused!\t:\t[%s]"
	LIKE_ALBUM_MESSAGE_TEMPLATE_LIKE_ALBUM_SUCCEEDED      = "✅ Like %s by [%s] succeeded!\t:\t[%s]"
)
