package constant

const (
	SEARCH_USE           = "search"
	SEARCH_HELP_TEMPLATE = `🔍 Search for the ID of content in Spotify.

You can set the args or the flag "query" to specify the search query.
If you set both args and flag "query", they will be combined.

You can set the flag "number" to limiting the number of search results.
Default is 1.

You can set the flag "type" to search type of the content.
If you don't set the flag "type", searching without specifying the content type will be executed.
You must specify the the flag "type" below :

  * 🎤 artist
  * 💿 album
  * 🎵 track

Usage:
  spotlike search [flags]

Flags:
  -h, --help           help for search
  -n, --number int     number of search results (default 1)
  -q, --query string   query for search
  -t, --type string    type of the content for search
`
	SEARCH_FLAG_QUERY              = "query"
	SEARCH_FLAG_QUERY_SHORTHAND    = "q"
	SEARCH_FLAG_QUERY_DEFAULT      = ""
	SEARCH_FLAG_QUERY_DESCRIPTION  = "query for search"
	SEARCH_FLAG_NUMBER             = "number"
	SEARCH_FLAG_NUMBER_SHORTHAND   = "n"
	SEARCH_FLAG_NUMBER_DEFAULT     = 1
	SEARCH_FLAG_NUMBER_DESCRIPTION = "number of search results"
	SEARCH_FLAG_TYPE               = "type"
	SEARCH_FLAG_TYPE_SHORTHAND     = "t"
	SEARCH_FLAG_TYPE_DEFAULT       = ""
	SEARCH_FLAG_TYPE_DESCRIPTION   = "type of the content for search"

	SEARCH_ERROR_MESSAGE_ARGS_OR_FLAG_QUERY_REQUIRED = `The arguments or the flag "query" is required...`
	SEARCH_ERROR_MESSAGE_FLAG_TYPE_INVALID           = `The argument of the flag "type" must be "artist", "album", or "track"...`
)
