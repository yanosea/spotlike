package track

import (
	"time"

	"github.com/zmb3/spotify/v2"
)

// Track is a struct that represents a Spotify track.
type Track struct {
	// ID is the Spotify ID of the track.
	ID spotify.ID
	// Name is the name of the track.
	Name string
	// Artists is a list of the artists that contributed to the track.
	Artists []spotify.SimpleArtist
	// Album is the album that the track belongs to.
	Album spotify.SimpleAlbum
	// TrackNumber is the track number of the track in the album.
	TrackNumber spotify.Numeric
	// ReleaseDate is the release date of the track.
	ReleaseDate time.Time
}

// NewTrack returns a new instance of Track struct.
func NewTrack(
	id spotify.ID,
	name string,
	artists []spotify.SimpleArtist,
	album spotify.SimpleAlbum,
	trackNumber spotify.Numeric,
	releaseDate time.Time,
) *Track {
	return &Track{
		ID:          id,
		Name:        name,
		Artists:     artists,
		Album:       album,
		TrackNumber: trackNumber,
		ReleaseDate: releaseDate,
	}
}
