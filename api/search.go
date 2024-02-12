package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

type SearchResult struct {
	Id          string
	Type        spotify.SearchType
	ArtistNames string
	AlbumName   string
	TrackName   string
}

const (
	search_error_message_something_wrong = "Something wrong occured..."
	search_error_message_not_found       = "The content [%s] was not found..."
)

func SearchByQuery(client *spotify.Client, searchType spotify.SearchType, query string) (*SearchResult, error) {
	result, err := client.Search(context.Background(), query, searchType, spotify.Limit(1))
	if err != nil {
		return nil, err
	}

	if result.Artists != nil {
		return &SearchResult{
			Id:          result.Artists.Artists[0].ID.String(),
			Type:        spotify.SearchTypeArtist,
			ArtistNames: result.Artists.Artists[0].Name,
		}, nil
	}

	if result.Albums != nil {
		return &SearchResult{
			Id:          result.Albums.Albums[0].ID.String(),
			Type:        spotify.SearchTypeAlbum,
			ArtistNames: combineArtistNames(result.Albums.Albums[0].Artists),
			AlbumName:   result.Albums.Albums[0].Name,
		}, nil
	}

	if result.Tracks != nil {
		return &SearchResult{
			Id:          result.Tracks.Tracks[0].ID.String(),
			Type:        spotify.SearchTypeTrack,
			ArtistNames: combineArtistNames(result.Tracks.Tracks[0].Artists),
			AlbumName:   result.Tracks.Tracks[0].Album.Name,
			TrackName:   result.Tracks.Tracks[0].Name,
		}, nil
	}

	return nil, errors.New(search_error_message_something_wrong)
}

func SearchById(client *spotify.Client, id string) (*SearchResult, error) {
	var unmarshalTypeErr *json.UnmarshalTypeError
	if _, err := client.GetArtist(context.Background(), spotify.ID(id)); errors.As(err, &unmarshalTypeErr) {
		// if search failed with UnmarshalTypeError, search artist albums again
		// c.f. https://github.com/zmb3/spotify/issues/243
		// c.f. https://community.spotify.com/t5/Spotify-for-Developers/Get-Artist-API-endpoint-responds-with-result-in-inconsistent/td-p/5806916
		if result, err := client.GetArtistAlbums(context.Background(), spotify.ID(id), nil); err == nil {
			for _, album := range result.Albums {
				for _, artist := range album.Artists {
					if artist.ID.String() == id {
						return &SearchResult{
							Id:          artist.ID.String(),
							Type:        spotify.SearchTypeArtist,
							ArtistNames: artist.Name,
						}, nil
					}
				}
			}
		}
	}

	if result, err := client.GetAlbum(context.Background(), spotify.ID(id)); err == nil {
		return &SearchResult{
			Id:          result.ID.String(),
			Type:        spotify.SearchTypeAlbum,
			ArtistNames: combineArtistNames(result.Artists),
			AlbumName:   result.Name,
		}, nil
	}

	if result, err := client.GetTrack(context.Background(), spotify.ID(id)); err == nil {
		return &SearchResult{
			Id:          result.ID.String(),
			Type:        spotify.SearchTypeTrack,
			ArtistNames: combineArtistNames(result.Artists),
			AlbumName:   result.Album.Name,
			TrackName:   result.Name,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf(search_error_message_not_found, id))
}
