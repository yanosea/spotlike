package constant

const (
	LIKE_ARTIST_USE           = "artist"
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
	LIKE_ARTIST_FLAG_ID                                    = "id"
	LIKE_ARTIST_SHORTHAND_ID                               = "i"
	LIKE_ARTIST_FLAG_DESCRIPTION_ID                        = "ID of the artist(s) for like"
	LIKE_ARTIST_FLAG_FORCE                                 = "force"
	LIKE_ARTIST_SHORTHAND_FORCE                            = "f"
	LIKE_ARTIST_FLAG_DESCRIPTION_FORCE                     = "like artist(s) without confirming"
	LIKE_ARTIST_FLAG_VERBOSE                               = "verbose"
	LIKE_ARTIST_SHORTHAND_VERBOSE                          = "v"
	LIKE_ARTIST_FLAG_DESCRIPTION_VERBOSE                   = "print verbose output"
	LIKE_ARTIST_ERROR_MESSAGE_TEMPLATE_ID_NOT_ARTIST       = "The ID you passed [%s] is not artist..."
	LIKE_ARTIST_MESSAGE_TEMPLATE_LIKE_ARTIST_ALREADY_LIKED = "✨ %s already liked!\t:\t[%s]"
	LIKE_ARTIST_MESSAGE_TEMPLATE_LIKE_ARTIST_REFUSED       = "❌ Like %s refused!\t:\t[%s]"
	LIKE_ARTIST_MESSAGE_TEMPLATE_LIKE_ARTIST_SUCCEEDED     = "✅ Like %s succeeded!\t:\t[%s]"
)
