package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/yanosea/spotlike/constants"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// SearchResult represents the search result from the Spotify API.
type SearchResult struct {
	// ID is the content's id
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

// SearchByQuery returns the search result by query.
func SearchByQuery(client *spotify.Client, searchType spotify.SearchType, query string) *SearchResult {
	// execute search
	result, err := client.Search(context.Background(), query, searchType, spotify.Limit(1))
	if err != nil {
		// search failed
		return &SearchResult{
			Result: false,
			Error:  err,
		}
	}

	if result.Artists != nil {
		// the type of the content is artist
		return &SearchResult{
			ID:          result.Artists.Artists[0].ID.String(),
			Type:        constants.Artist,
			ArtistNames: result.Artists.Artists[0].Name,
			Result:      true,
		}
	} else if result.Albums != nil {
		// the type of the content is album
		return &SearchResult{
			ID:          result.Albums.Albums[0].ID.String(),
			Type:        constants.Album,
			ArtistNames: combineArtistNames(result.Albums.Albums[0].Artists),
			AlbumName:   result.Albums.Albums[0].Name,
			Result:      true,
		}
	} else if result.Tracks != nil {
		// the type of the content is track
		return &SearchResult{
			ID:          result.Tracks.Tracks[0].ID.String(),
			Type:        constants.Track,
			ArtistNames: combineArtistNames(result.Tracks.Tracks[0].Artists),
			AlbumName:   result.Tracks.Tracks[0].Album.Name,
			TrackName:   result.Tracks.Tracks[0].Name,
			Result:      true,
			Error:       nil,
		}
	} else {
		// search failed
		return &SearchResult{
			Result: false,
			Error:  errors.New(constants.SearchFailedErrorMessage),
		}
	}
}

// SearchById returns the search result by ID.
func SearchById(client *spotify.Client, id string) *SearchResult {
	// execute the search
	if result, err := client.GetArtist(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is artist
		return &SearchResult{
			ID:          result.ID.String(),
			Type:        constants.Artist,
			ArtistNames: result.Name,
			Result:      true,
		}
	} else if result, err := client.GetAlbum(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is album
		return &SearchResult{
			ID:          result.ID.String(),
			Type:        constants.Album,
			ArtistNames: combineArtistNames(result.Artists),
			AlbumName:   result.Name,
			Result:      true,
		}
	} else if result, err := client.GetTrack(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is track
		return &SearchResult{
			ID:          result.ID.String(),
			Type:        constants.Track,
			ArtistNames: combineArtistNames(result.Artists),
			AlbumName:   result.Album.Name,
			TrackName:   result.Name,
			Result:      true,
		}
	} else {
		// content not found
		return &SearchResult{
			Result: false,
			Error:  errors.New(fmt.Sprintf(constants.SearchFailedNotFoundErrorMessageFormat, id)),
		}
	}
}
