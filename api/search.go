package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/yanosea/spotlike/app"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// SearchResult represents the search result from the Spotify API.
type SearchResult struct {
	// ID is content's id
	ID string
	// Type is the content type
	Type string
	// ArtistNames is the names of the artists
	ArtistNames string
	// AlbumName is the name of the album
	AlbumName string
	// TrackName is the name of the track
	TrackName string
	// Result is the search result (true: succeeded, false: failed)
	Result bool
	// Error is the error returned from the Spotify API
	Error error
}

// searchResult holds the current search result.
var searchResult *SearchResult

// SearchByQuery returns the search result by query.
func SearchByQuery(client *spotify.Client, searchType spotify.SearchType, query string) *SearchResult {
	// execute search
	if result, err := client.Search(context.Background(), query, searchType, spotify.Limit(1)); err != nil {
		// search failed
		searchResult = &SearchResult{
			Result: false,
			Error:  err,
		}
	} else {
		if result.Artists != nil {
			// the type of the content is artist
			searchResult = &SearchResult{
				ID:          result.Artists.Artists[0].ID.String(),
				Type:        "Artist",
				ArtistNames: result.Artists.Artists[0].Name,
				Result:      true,
			}
		} else if result.Albums != nil {
			// the type of the content is album
			searchResult = &SearchResult{
				ID:          result.Albums.Albums[0].ID.String(),
				Type:        "Album",
				ArtistNames: app.CombineArtistNames(result.Albums.Albums[0].Artists),
				AlbumName:   result.Albums.Albums[0].Name,
				Result:      true,
			}
		} else if result.Tracks != nil {

			// the type of the content is track
			searchResult = &SearchResult{
				ID:          result.Tracks.Tracks[0].ID.String(),
				Type:        "Track",
				ArtistNames: app.CombineArtistNames(result.Tracks.Tracks[0].Artists),
				AlbumName:   result.Tracks.Tracks[0].Album.Name,
				TrackName:   result.Tracks.Tracks[0].Name,
				Result:      true,
				Error:       nil,
			}
		} else {
			// search failed
			searchResult = &SearchResult{
				Result: false,
				Error:  errors.New("Search result is wrong."),
			}
		}
	}
	return searchResult
}

// SearchById returns the search result by ID.
func SearchById(client *spotify.Client, id string) *SearchResult {
	// execute the search
	if result, err := client.GetArtist(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is artist
		searchResult = &SearchResult{
			ID:          result.ID.String(),
			Type:        "Artist",
			ArtistNames: result.Name,
			Result:      true,
		}
	} else if result, err := client.GetAlbum(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is album
		searchResult = &SearchResult{
			ID:          result.ID.String(),
			Type:        "Album",
			ArtistNames: app.CombineArtistNames(result.Artists),
			AlbumName:   result.Name,
			Result:      true,
		}
	} else if result, err := client.GetTrack(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is track
		searchResult = &SearchResult{
			ID:          result.ID.String(),
			Type:        "Track",
			ArtistNames: app.CombineArtistNames(result.Artists),
			AlbumName:   result.Album.Name,
			TrackName:   result.Name,
			Result:      true,
		}
	} else {
		// content not found
		searchResult = &SearchResult{
			Result: false,
			Error:  errors.New(fmt.Sprintf("The content [%s] was not found", id)),
		}
	}
	return searchResult
}
