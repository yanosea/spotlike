package album

import (
	"time"

	"github.com/zmb3/spotify/v2"
)

// Album is a struct that represents a Spotify album.
type Album struct {
	// ID is the Spotify ID of the album.
	ID spotify.ID
	// Name is the name of the album.
	Name string
	// Artists is a list of the artists that contributed to the album.
	Artists []spotify.SimpleArtist
	// ReleaseDate is the release date of the album.
	ReleaseDate time.Time
}

// NewAlbum returns a new instance of Album struct.
func NewAlbum(
	id spotify.ID,
	name string,
	artists []spotify.SimpleArtist,
	releaseDate time.Time,
) *Album {
	return &Album{
		ID:          id,
		Name:        name,
		Artists:     artists,
		ReleaseDate: releaseDate,
	}
}
