package util

import (
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// constants
const (
	// STRING_ID is the STRING of ID
	STRING_ID = "ID"
	// STRING_TYPE is the STRING of Type
	STRING_TYPE = "Type"
	// STRING_ARTIST is the STRING of Artist
	STRING_ARTIST = "Artist"
	// STRING_ALBUM is the STRING of Album
	STRING_ALBUM = "Album"
	// STRING_TRACK is the STRING of Track
	STRING_TRACK = "Track"
)

// valiables
var (
	// SEARCH_TYPE_MAP maps string values to spotify.SearchType values.
	SEARCH_TYPE_MAP = map[string]spotify.SearchType{
		"artist": spotify.SearchTypeArtist,
		"album":  spotify.SearchTypeAlbum,
		"track":  spotify.SearchTypeTrack,
	}
	// SEARCH_TYPE_MAP_REVERSED maps spotify.SearchType to string values.
	SEARCH_TYPE_MAP_REVERSED = map[spotify.SearchType]string{
		spotify.SearchTypeArtist: "artist",
		spotify.SearchTypeAlbum:  "album",
		spotify.SearchTypeTrack:  "track",
	}
)
