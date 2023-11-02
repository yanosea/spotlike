/*
Package api provides functions for sending requests to the Spotify API
*/
package api

import (
	"context"
	"fmt"
	"sort"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// LikeResult represents the result of liking content using the Spotify API.
type LikeResult struct {
	// ID is content's id
	ID string
	// Type is the content type
	Type string
	// ArtistNames is the names of the artists
	ArtistNames string
	// AlbumName is the album name
	AlbumName string
	// TrackName is the track name
	TrackName string
	// Result is the like result (true: succeeded, false: failed)
	Result bool
	// Error is the error returned from the Spotify API
	Error error
	// ErrorMessage is the Error message for the like result
	ErrorMessage string
}

// TrackWithAlbumName represents a Spotify simple track with an album name.
type TrackWithAlbumName struct {
	// Track is Spotify simple track
	Track spotify.SimpleTrack
	// AlbumName is the Name of the album the track is included in
	AlbumName string
}

// likeResults holds the like results.
var likeResults []*LikeResult

// LikeArtistById returns the like result for an artist with the given ID.
func LikeArtistById(client *spotify.Client, id string) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Artist" {
		// execute like
		if err := client.FollowArtist(context.Background(), spotify.ID(sr.ID)); err != nil {
			// like failed
			likeResults = append(likeResults, &LikeResult{
				ID:          sr.ID,
				Type:        "Artist",
				ArtistNames: sr.ArtistNames,
				Result:      false,
				Error:       err,
			})
		} else {
			// like succeeded
			likeResults = append(likeResults, &LikeResult{
				ID:          sr.ID,
				Type:        "Artist",
				ArtistNames: sr.ArtistNames,
				Result:      true,
			})
		}
	} else {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the artist by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeAllAlbumsReleasedByArtistById returns the like results for all albums released by an artist with the given ID.
func LikeAllAlbumsReleasedByArtistById(client *spotify.Client, id string) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Artist" {
		// get all albums released by the artist
		if allAlbums, err := client.GetArtistAlbums(context.Background(), spotify.ID(sr.ID), nil); err != nil {
			likeResults = append(likeResults, &LikeResult{
				// getting all albums by the artist searched by ID failed
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Get all albums by the artist searched by ID failed...\t:\t[%s]", id),
			})
		} else {
			// sort albums by release date
			sort.Slice(allAlbums.Albums, func(i, j int) bool {
				return allAlbums.Albums[i].ReleaseDateTime().Before(allAlbums.Albums[j].ReleaseDateTime())
			})
			// execute like
			for _, album := range allAlbums.Albums {
				if err := client.AddAlbumsToLibrary(context.Background(), album.ID); err != nil {
					likeResults = append(likeResults, &LikeResult{
						// like failed
						ID:          album.ID.String(),
						Type:        "Album",
						ArtistNames: sr.ArtistNames,
						AlbumName:   album.Name,
						Result:      false,
						Error:       err,
					})
				} else {
					// like succeeded
					likeResults = append(likeResults, &LikeResult{
						ID:          album.ID.String(),
						Type:        "Album",
						ArtistNames: sr.ArtistNames,
						AlbumName:   album.Name,
						Result:      true,
					})
				}
			}
		}
	} else {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the artist by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeAllTracksReleasedByArtistById returns the like results for all tracks released by an artist with the given ID.
func LikeAllTracksReleasedByArtistById(client *spotify.Client, id string) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Artist" {
		// get all albums released by the artist
		if allAlbums, err := client.GetArtistAlbums(context.Background(), spotify.ID(sr.ID), nil); err != nil {
			likeResults = append(likeResults, &LikeResult{
				// getting all albums by the artist searched by ID failed
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Get all albums by the artist searched by ID failed...\t:\t[%s]", id),
			})
		} else {
			// sort albums by release date
			sort.Slice(allAlbums.Albums, func(i, j int) bool {
				return allAlbums.Albums[i].ReleaseDateTime().Before(allAlbums.Albums[j].ReleaseDateTime())
			})

			// get all tracks from all albums
			var allTracks []TrackWithAlbumName
			for _, album := range allAlbums.Albums {
				if tracks, err := client.GetAlbumTracks(context.Background(), album.ID); err != nil {
					// getting all tracks in all albums by the artist searched by ID failed
					likeResults = append(likeResults, &LikeResult{
						Result:       false,
						Error:        err,
						ErrorMessage: fmt.Sprintf("Get all tracks in all albums by the artist searched by ID failed...\t:\t[%s]", id),
					})
				} else {
					for _, track := range tracks.Tracks {
						trackWithAlbumName := &TrackWithAlbumName{
							Track:     track,
							AlbumName: album.Name,
						}
						allTracks = append(allTracks, *trackWithAlbumName)
					}
				}
			}

			// execute like
			for _, track := range allTracks {
				if err := client.AddTracksToLibrary(context.Background(), track.Track.ID); err != nil {
					// like failed
					likeResults = append(likeResults, &LikeResult{
						ID:          track.Track.ID.String(),
						Type:        "Track",
						ArtistNames: sr.ArtistNames,
						AlbumName:   track.AlbumName,
						TrackName:   track.Track.Name,
						Result:      false,
						Error:       err,
					})
				} else {
					// like succeeded
					likeResults = append(likeResults, &LikeResult{
						ID:          track.Track.ID.String(),
						Type:        "Track",
						ArtistNames: sr.ArtistNames,
						AlbumName:   track.AlbumName,
						TrackName:   track.Track.Name,
						Result:      true,
					})
				}
			}
		}
	} else {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the artist by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeAlbumById returns an error if liking an album with the given ID is failed.
func LikeAlbumById(client *spotify.Client, id string) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Album" {
		// execute like
		if err := client.AddAlbumsToLibrary(context.Background(), spotify.ID(id)); err != nil {
			// like failed
			likeResults = append(likeResults, &LikeResult{
				ID:          id,
				Type:        "Album",
				ArtistNames: sr.ArtistNames,
				AlbumName:   sr.AlbumName,
				Result:      false,
				Error:       err,
			})
		} else {
			// like failed
			likeResults = append(likeResults, &LikeResult{
				ID:          id,
				Type:        "Album",
				ArtistNames: sr.ArtistNames,
				AlbumName:   sr.AlbumName,
				Result:      true,
			})
		}
	} else {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the album by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeAllTracksInAlbumById returns an error if liking all tracks in an album with the given ID is failed.
func LikeAllTracksInAlbumById(client *spotify.Client, id string) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Album" {
		// get all tracks in the album
		if allTracks, err := client.GetAlbumTracks(context.Background(), spotify.ID(id)); err != nil {
			// getting all tracks in the album searched by ID failed
			likeResults = append(likeResults, &LikeResult{
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Get all tracks in the album searched by ID failed...\t:\t[%s]", id),
			})
		} else {
			for _, track := range allTracks.Tracks {
				// execute like
				if err := client.AddTracksToLibrary(context.Background(), track.ID); err != nil {
					// like failed
					likeResults = append(likeResults, &LikeResult{
						ID:          track.ID.String(),
						Type:        "Track",
						ArtistNames: sr.ArtistNames,
						AlbumName:   sr.AlbumName,
						TrackName:   track.Name,
						Result:      false,
						Error:       err,
					})
				} else {
					// like succeeded
					likeResults = append(likeResults, &LikeResult{
						ID:          track.ID.String(),
						Type:        "Track",
						ArtistNames: sr.ArtistNames,
						AlbumName:   sr.AlbumName,
						TrackName:   track.Name,
						Result:      true,
					})
				}
			}
		}
	} else {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the album by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeTrackById returns an error if liking a track with the given ID is failed.
func LikeTrackById(client *spotify.Client, id string) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Track" {
		// execute like
		if err := client.AddTracksToLibrary(context.Background(), spotify.ID(id)); err != nil {
			likeResults = append(likeResults, &LikeResult{
				// like failed
				ID:          sr.ID,
				Type:        "Track",
				ArtistNames: sr.ArtistNames,
				AlbumName:   sr.AlbumName,
				TrackName:   sr.TrackName,
				Result:      false,
				Error:       err,
			})
		} else {
			// like succeeded
			likeResults = append(likeResults, &LikeResult{
				ID:          sr.ID,
				Type:        "Track",
				ArtistNames: sr.ArtistNames,
				AlbumName:   sr.AlbumName,
				TrackName:   sr.TrackName,
				Result:      true,
			})
		}
	} else {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the track by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}
