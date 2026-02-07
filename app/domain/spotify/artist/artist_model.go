package artist

import (
	"github.com/zmb3/spotify/v2"
)

// Artist is a struct that represents a Spotify artist.
type Artist struct {
	// ID is the Spotify ID of the artist.
	ID spotify.ID
	// Name is the name of the artist.
	Name string
}

// NewArtist returns a new instance of Artist struct.
func NewArtist(
	id spotify.ID,
	name string,
) *Artist {
	return &Artist{
		ID:   id,
		Name: name,
	}
}
