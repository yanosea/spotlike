package api

import (
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

func combineArtistNames(artists []spotify.SimpleArtist) string {
	var artistNames string
	for index, artist := range artists {
		artistNames += artist.Name
		if index+1 != len(artists) {
			artistNames += ", "
		}
	}

	return artistNames
}
