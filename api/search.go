package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// SearchResult represents the search result from the Spotify API.
type SearchResult struct {
	// Id is the content's id
	Id string
	// Type is the content type
	Type spotify.SearchType
	// ArtistNames is the names of the artists
	ArtistNames string
	// AlbumName is the name of the album
	AlbumName string
	// TrackName is the name of the track
	TrackName string
}

// constants
const (
	// search_error_message_something_wrong is the error message for something wrong searching.
	search_error_message_something_wrong = "Something wrong occured..."
	// search_error_message_not_found is the error message for not found.
	search_error_message_not_found = "The content [%s] was not found..."
)

// SearchByQuery returns the search result by query.
func SearchByQuery(client *spotify.Client, searchType spotify.SearchType, query string) (*SearchResult, error) {
	// execute search
	result, err := client.Search(context.Background(), query, searchType, spotify.Limit(1))
	if err != nil {
		// search failed
		return nil, err
	}

	if result.Artists != nil {
		// the type of the content is artist
		return &SearchResult{
			Id:          result.Artists.Artists[0].ID.String(),
			Type:        spotify.SearchTypeArtist,
			ArtistNames: result.Artists.Artists[0].Name,
		}, nil
	} else if result.Albums != nil {
		// the type of the content is album
		return &SearchResult{
			Id:          result.Albums.Albums[0].ID.String(),
			Type:        spotify.SearchTypeAlbum,
			ArtistNames: combineArtistNames(result.Albums.Albums[0].Artists),
			AlbumName:   result.Albums.Albums[0].Name,
		}, nil
	} else if result.Tracks != nil {
		// the type of the content is track
		return &SearchResult{
			Id:          result.Tracks.Tracks[0].ID.String(),
			Type:        spotify.SearchTypeTrack,
			ArtistNames: combineArtistNames(result.Tracks.Tracks[0].Artists),
			AlbumName:   result.Tracks.Tracks[0].Album.Name,
			TrackName:   result.Tracks.Tracks[0].Name,
		}, nil
	} else {
		// search failed
		return nil, errors.New(search_error_message_something_wrong)
	}
}

// SearchById returns the search result by ID.
func SearchById(client *spotify.Client, id string) (*SearchResult, error) {
	var unmarshalTypeErr *json.UnmarshalTypeError
	// execute the search
	if result, err := client.GetArtist(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is artist
		return &SearchResult{
			Id:          result.ID.String(),
			Type:        spotify.SearchTypeArtist,
			ArtistNames: result.Name,
		}, nil
	} else if errors.As(err, &unmarshalTypeErr) {
		// if search failed with UnmarshalTypeError, search artist albums again
		if result, err := client.GetArtistAlbums(context.Background(), spotify.ID(id), nil); err == nil {
			// search for an album with a single artist from the retrieved albums
			for _, album := range result.Albums {
				if len(album.Artists) == 1 {
					// the type of the content is artist
					return &SearchResult{
						Id:          album.Artists[0].ID.String(),
						Type:        spotify.SearchTypeArtist,
						ArtistNames: album.Artists[0].Name,
					}, nil
				}
			}
		}
	}

	if result, err := client.GetAlbum(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is album
		return &SearchResult{
			Id:          result.ID.String(),
			Type:        spotify.SearchTypeAlbum,
			ArtistNames: combineArtistNames(result.Artists),
			AlbumName:   result.Name,
		}, nil
	} else if result, err := client.GetTrack(context.Background(), spotify.ID(id)); err == nil {
		// the type of the content is track
		return &SearchResult{
			Id:          result.ID.String(),
			Type:        spotify.SearchTypeTrack,
			ArtistNames: combineArtistNames(result.Artists),
			AlbumName:   result.Album.Name,
			TrackName:   result.Name,
		}, nil
	} else {
		// content not found
		return nil, errors.New(fmt.Sprintf(search_error_message_not_found, id))
	}
}
