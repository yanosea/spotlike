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
	ReleaseDate string
	AlbumName   string
	TrackName   string
}

const (
	search_error_message_not_artist_album_track = "The content was not artist, album, or track..."
	search_error_message_template_id_not_found  = "The content [%s] was not found..."
)

func SearchByQuery(client *spotify.Client, query string, number int, searchType spotify.SearchType) ([]SearchResult, error) {
	var searchResultList []SearchResult
	result, err := client.Search(context.Background(), query, searchType, spotify.Limit(number))
	if err != nil {
		return nil, err
	}

	if result.Artists != nil {
		for _, artist := range result.Artists.Artists {
			searchResultList = append(searchResultList, SearchResult{
				Id:          artist.ID.String(),
				Type:        spotify.SearchTypeArtist,
				ArtistNames: artist.Name,
			})
		}
		return searchResultList, nil
	}

	if result.Albums != nil {
		for _, album := range result.Albums.Albums {
			searchResultList = append(searchResultList, SearchResult{
				Id:          album.ID.String(),
				Type:        spotify.SearchTypeAlbum,
				ArtistNames: combineArtistNames(album.Artists),
				ReleaseDate: album.ReleaseDate,
				AlbumName:   album.Name,
			})
		}
		return searchResultList, nil
	}

	if result.Tracks != nil {
		for _, track := range result.Tracks.Tracks {
			searchResultList = append(searchResultList, SearchResult{
				Id:          track.ID.String(),
				Type:        spotify.SearchTypeTrack,
				ArtistNames: combineArtistNames(track.Artists),
				ReleaseDate: track.Album.ReleaseDate,
				AlbumName:   track.Album.Name,
				TrackName:   track.Name,
			})
		}
		return searchResultList, nil
	}

	return nil, errors.New(search_error_message_not_artist_album_track)
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
			ReleaseDate: result.ReleaseDate,
			AlbumName:   result.Name,
		}, nil
	}

	if result, err := client.GetTrack(context.Background(), spotify.ID(id)); err == nil {
		return &SearchResult{
			Id:          result.ID.String(),
			Type:        spotify.SearchTypeTrack,
			ArtistNames: combineArtistNames(result.Artists),
			ReleaseDate: result.Album.ReleaseDate,
			AlbumName:   result.Album.Name,
			TrackName:   result.Name,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf(search_error_message_template_id_not_found, id))
}
