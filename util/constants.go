package util

import (
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

const (
	STRING_ID      = "ID"
	STRING_TYPE    = "Type"
	STRING_ARTIST  = "Artist"
	STRING_ALBUM   = "Album"
	STRING_RELEASE = "Release"
	STRING_TRACK   = "Track"

	AUTH_ENV_SPOTIFY_ID            = "SPOTIFY_ID"
	AUTH_ENV_SPOTIFY_SECRET        = "SPOTIFY_SECRET"
	AUTH_ENV_SPOTIFY_REDIRECT_URI  = "SPOTIFY_REDIRECT_URI"
	AUTH_ENV_SPOTIFY_REFRESH_TOKEN = "SPOTIFY_REFRESH_TOKEN"
)

var (
	SEARCH_TYPE_MAP = map[string]spotify.SearchType{
		"artist": spotify.SearchTypeArtist,
		"album":  spotify.SearchTypeAlbum,
		"track":  spotify.SearchTypeTrack,
	}
	SEARCH_TYPE_MAP_REVERSED = map[spotify.SearchType]string{
		spotify.SearchTypeArtist: "artist",
		spotify.SearchTypeAlbum:  "album",
		spotify.SearchTypeTrack:  "track",
	}
)
