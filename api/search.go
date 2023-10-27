/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package api

import (
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// SearchResult : search result from spotify api
type SearchResult struct {
	// ID : content's id
	ID string
	// Type : content type
	Type string
	// Name : content name
	Name string
	// Album : name of album on which track is inclued
	Album string
	// Artist : artist of album or track
	Artist string
}

// Search : returns search with spotify library
func Search(spt *SpotifyClient, searchType spotify.SearchType, query string) (*SearchResult, error) {
	// execute search
	if result, err := spt.Client.Search(spt.Context, query, searchType, spotify.Limit(1)); err != nil {
		return nil, err
	} else {
		var searchResult *SearchResult
		// artist
		if result.Artists != nil {
			searchResult = &SearchResult{
				ID:   result.Artists.Artists[0].ID.String(),
				Type: "Artist",
				Name: result.Artists.Artists[0].Name,
			}
		}

		// album
		if result.Albums != nil {
			searchResult = &SearchResult{
				ID:     result.Albums.Albums[0].ID.String(),
				Type:   "Album",
				Name:   result.Albums.Albums[0].Name,
				Artist: result.Albums.Albums[0].Artists[0].Name,
			}
		}

		// track
		if result.Tracks != nil {
			searchResult = &SearchResult{
				ID:     result.Tracks.Tracks[0].ID.String(),
				Type:   "Track",
				Name:   result.Tracks.Tracks[0].Name,
				Album:  result.Tracks.Tracks[0].Album.Name,
				Artist: result.Tracks.Tracks[0].Artists[0].Name,
			}
		}

		return searchResult, nil
	}
}
