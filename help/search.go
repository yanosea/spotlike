package help

const (
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
)
