package constants

import (
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// constants
const (
	// EnvSpotifyID is env string of the Spotify client ID
	EnvSpotifyID = "SPOTIFY_ID"
	// EnvSpotifyIDInputLabel is the label of the Spotify client ID
	EnvSpotifyIDInputLabel = "Input your Spotify Client ID"
	// EnvSpotifySecret is env string of the Spotify client secret
	EnvSpotifySecret = "SPOTIFY_SECRET"
	// EnvSpotifySecretInputLabel is the label of the Spotify client secret
	EnvSpotifySecretInputLabel = "Input your Spotify Client ID"
	// EnvSpotifySecretMaskCharacter is the character for mask the Spotify client secret
	EnvSpotifySecretMaskCharacter = '*'
	// EnvSpotifyRedirectUri is env string of the Spotify redirect URI
	EnvSpotifyRedirectUri = "SPOTIFY_REDIRECT_URI"
	// EnvSpotifySecretInputLabel is the label of the Spotify redirect URI
	EnvSpotifyRedirectUriInputLabel = "Input your Spotify Redirect URI"
	// EnvSpotifyRefreshToken is env string of the Spotify refresh token
	EnvSpotifyRefreshToken = "SPOTIFY_REFRESH_TOKEN"
	// EnvSpotifyRefreshTokenInputLabel is the label of the Spotify refresh token
	EnvSpotifyRefreshTokenInputLabel = "Input your Spotify Refresh Token if you have one (if you don't have one, leave it empty and press enter.)"

	// Id is the string of ID
	Id = "ID"
	// Type is the string of Type
	Type = "Type"
	// Artist is the string of Artist
	Artist = "Artist"
	// Album is the string of Album
	Album = "Album"
	// Track is the string of Track
	Track = "Track"

	// RootUse is the useage of root command.
	RootUse = "spotlike"
	// RootShort is the short description of root command.
	RootShort = "spotlike is the CLI tool to LIKE contents in Spotify."
	// long is the long description of root command.
	RootLong = `'spotlike' is the CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`

	// SearchUse is the useage of search command.
	SearchUse = "search"
	// SearchShort is the short description of search command.
	SearchShort = "Search for the ID of content in Spotify."
	// SearchLong is the long description of search command.
	SearchLong = `Search for the ID of content in Spotify.

You can search for content using the type option below:
  * artist
  * album
  * track`
	// SearchFlagType is the string of the type flag of the search command.
	SearchFlagType = "type"
	// SearchFlagTypeShorthand  is the string of the type shorthand flag of the search command.
	SearchFlagTypeShorthand = "t"
	// SearchFlagTypeDescription  is the description of the type flag of the search command.
	SearchFlagTypeDescription = "type of the content for search"
	// SearchFlagQuery is the string of the query flag of the search command.
	SearchFlagQuery = "query"
	// SearchFlagQueryShorthand is the string of the query shorthand flag of the search command.
	SearchFlagQueryShorthand = "q"
	// SearchFlagQueryDescription the description of the query flag of the search command.
	SearchFlagQueryDescription = "query for search"
	// SearchFlagTypeInvalidErrorMessage is the error message for the invalid type.
	SearchFlagTypeInvalidErrorMessage = `the argument of the flag "type" must be "artist", "album", or "track..."`
	// SearchFailedErrorMessageFormat     is the error message format for search failure.
	SearchFailedErrorResultMessageFormat   = "Search for %s failed..."
	SearchFailedErrorMessage               = "Search result is wrong..."
	SearchFailedNotFoundErrorMessageFormat = "The content [%s] was not found..."
)

// SearchTypeMap maps strings to spotify.SearchType values.
var SearchTypeMap = map[string]spotify.SearchType{
	"artist": spotify.SearchTypeArtist,
	"album":  spotify.SearchTypeAlbum,
	"track":  spotify.SearchTypeTrack,
}
