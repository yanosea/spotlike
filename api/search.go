/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package api

import (
	"errors"

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
	Result bool
	Error  error
}

// searchResult : search result from spotify api
var searchResult *SearchResult

// SearchByQuery : returns the search result by query
func SearchByQuery(spt *SpotifyClient, searchType spotify.SearchType, query string) *SearchResult {
	// execute search
	if result, err := spt.Client.Search(spt.Context, query, searchType, spotify.Limit(1)); err != nil {
		searchResult = &SearchResult{
			ID:     "",
			Type:   "Artist",
			Name:   "",
			Result: false,
			Error:  err,
		}
	} else {
		// artist
		if result.Artists != nil {
			searchResult = &SearchResult{
				ID:     result.Artists.Artists[0].ID.String(),
				Type:   "Artist",
				Name:   result.Artists.Artists[0].Name,
				Result: true,
				Error:  nil,
			}
		}

		// album
		if result.Albums != nil {
			searchResult = &SearchResult{
				ID:     result.Albums.Albums[0].ID.String(),
				Type:   "Album",
				Name:   result.Albums.Albums[0].Name,
				Artist: result.Albums.Albums[0].Artists[0].Name,
				Result: true,
				Error:  nil,
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
				Result: true,
				Error:  nil,
			}
		}
	}

	return searchResult
}

// SearchById : returns the search result by ID
func SearchById(spt *SpotifyClient, id string) *SearchResult {
	// execute search
	if result, err := spt.Client.GetArtist(spt.Context, spotify.ID(id)); err == nil {
		// artist
		searchResult = &SearchResult{
			ID:   result.ID.String(),
			Type: "Artist",
			Name: result.Name,
		}
	} else if result, err := spt.Client.GetAlbum(spt.Context, spotify.ID(id)); err == nil {
		// album
		searchResult = &SearchResult{
			ID:     result.ID.String(),
			Type:   "Album",
			Name:   result.Name,
			Artist: result.Artists[0].Name,
		}
	} else if result, err := spt.Client.GetTrack(spt.Context, spotify.ID(id)); err == nil {
		// track
		searchResult = &SearchResult{
			ID:     result.ID.String(),
			Type:   "Track",
			Name:   result.Name,
			Album:  result.Album.Name,
			Artist: result.Artists[0].Name,
		}
	} else {
		// content not found
		searchResult = &SearchResult{
			ID:     "",
			Type:   "",
			Name:   "",
			Result: false,
			Error:  errors.New("the content was not found"),
		}
	}

	return searchResult
}
