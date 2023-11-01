/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package api

import (
	"context"
	"fmt"
	"sort"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// LikeResult : like result from spotify api
type LikeResult struct {
	// ID : content's ID
	ID string
	// Type : content type
	Type string
	// ArtistName : artist name
	ArtistName string
	// AlbumName : album name
	AlbumName string
	// TrackName : track name
	TrackName string
	// Result : like result (true : succeeded / false : failed)
	Result bool
	// Error : error returned from spotify api
	Error error
	// ErrorMessage : error message for like result
	ErrorMessage string
}

// TrackWithAlbumName : spotify simple track with album name
type TrackWithAlbumName struct {
	// Track : spotify simple track
	Track spotify.SimpleTrack
	// AlbumName : name of album the track included
	AlbumName string
}

// likelikeResults : like results from spotify api
var likeResults []*LikeResult

// LikeArtistById : returns like result
func LikeArtistById(client *spotify.Client, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(client, id); sr.Result && sr.Type == "Artist" {
		// execute like artist
		if err := client.FollowArtist(context.Background(), spotify.ID(sr.ID)); err != nil {
			// like the artist was failed
			likeResults = append(likeResults, &LikeResult{
				ID:         sr.ID,
				Type:       "Artist",
				ArtistName: sr.Name,
				Result:     false,
				Error:      err,
			})
		} else {
			// like the artist was succeeded
			likeResults = append(likeResults, &LikeResult{
				ID:         sr.ID,
				Type:       "Artist",
				ArtistName: sr.Name,
				Result:     true,
			})
		}
	} else {
		// search the artist by id was failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the artist by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeAllAlbumsReleasedByArtistById : returns like results
func LikeAllAlbumsReleasedByArtistById(client *spotify.Client, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(client, id); sr.Result && sr.Type == "Artist" {
		// get all albums released by the artist
		if allAlbums, err := client.GetArtistAlbums(context.Background(), spotify.ID(sr.ID), nil); err != nil {
			likeResults = append(likeResults, &LikeResult{
				// get all albums by the artist searched by id was failed
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Get all albums by the artist searched by ID failed...\t:\t[%s]", id),
			})
		} else {
			// sort albums by release date
			sort.Slice(allAlbums.Albums, func(i, j int) bool {
				return allAlbums.Albums[i].ReleaseDateTime().Before(allAlbums.Albums[j].ReleaseDateTime())
			})
			// like by album
			for _, album := range allAlbums.Albums {
				if client.AddAlbumsToLibrary(context.Background(), album.ID); err != nil {
					likeResults = append(likeResults, &LikeResult{
						// like the album by the artist searched by id was failed
						ID:         album.ID.String(),
						Type:       "Album",
						ArtistName: sr.Name,
						AlbumName:  album.Name,
						Result:     false,
						Error:      err,
					})
				} else {
					// like the album by the artist searched by id was succeeded
					likeResults = append(likeResults, &LikeResult{
						ID:         album.ID.String(),
						Type:       "Album",
						ArtistName: sr.Name,
						AlbumName:  album.Name,
						Result:     true,
					})
				}
			}
		}
	} else {
		// search the artist by id was failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the artist by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeAllTracksReleasedByArtistById : returns results
func LikeAllTracksReleasedByArtistById(client *spotify.Client, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(client, id); sr.Result && sr.Type == "Artist" {
		// get all albums released by the artist
		if allAlbums, err := client.GetArtistAlbums(context.Background(), spotify.ID(sr.ID), nil); err != nil {
			// get all albums by the artist searched by id was failed
			likeResults = append(likeResults, &LikeResult{
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
					// get all tracks in all albums by the artist searched by id was failed
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

			// like by track
			for _, track := range allTracks {
				if client.AddTracksToLibrary(context.Background(), track.Track.ID); err != nil {
					// like the track in all albums by the artist searched by id was failed
					likeResults = append(likeResults, &LikeResult{
						ID:         track.Track.ID.String(),
						Type:       "Track",
						ArtistName: sr.Name,
						AlbumName:  track.AlbumName,
						TrackName:  track.Track.Name,
						Result:     false,
						Error:      err,
					})
				} else {
					// like the track in all albums by the artist searched by id was succeeded
					likeResults = append(likeResults, &LikeResult{
						ID:         track.Track.ID.String(),
						Type:       "Track",
						ArtistName: sr.Name,
						AlbumName:  track.AlbumName,
						TrackName:  track.Track.Name,
						Result:     true,
					})
				}
			}
		}
	} else {
		// search the artist by id was failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the artist by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeAlbumById : returns error if liking is failed
func LikeAlbumById(client *spotify.Client, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(client, id); sr.Result && sr.Type == "Album" {
		// execute like album
		if err := client.AddAlbumsToLibrary(context.Background(), spotify.ID(id)); err != nil {
			// like the album was failed
			likeResults = append(likeResults, &LikeResult{
				ID:         id,
				Type:       "Album",
				ArtistName: sr.Artist,
				AlbumName:  sr.Name,
				Result:     false,
				Error:      err,
			})
		} else {
			// like the album was succeeded
			likeResults = append(likeResults, &LikeResult{
				ID:         id,
				Type:       "Album",
				ArtistName: sr.Artist,
				AlbumName:  sr.Name,
				Result:     true,
			})
		}
	} else {
		// search the album by id was failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the album by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeAllTracksInAlbumById : returns error if liking is failed
func LikeAllTracksInAlbumById(client *spotify.Client, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(client, id); sr.Result && sr.Type == "Album" {
		// get all tracks on the artist
		if allTracks, err := client.GetAlbumTracks(context.Background(), spotify.ID(id)); err != nil {
			// get all tracks in the album searched by id was failed
			likeResults = append(likeResults, &LikeResult{
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Get all tracks in the album searched by ID failed...\t:\t[%s]", id),
			})
		} else {
			for _, track := range allTracks.Tracks {
				// execute like track
				if client.AddTracksToLibrary(context.Background(), track.ID); err != nil {
					// like the track in the album searched by id was failed
					likeResults = append(likeResults, &LikeResult{
						ID:         track.ID.String(),
						Type:       "Track",
						ArtistName: sr.Name,
						AlbumName:  sr.Name,
						TrackName:  track.Name,
						Result:     false,
						Error:      err,
					})
				} else {
					// like the track in the album searched by id was succeeded
					likeResults = append(likeResults, &LikeResult{
						ID:         track.ID.String(),
						Type:       "Track",
						ArtistName: sr.Name,
						AlbumName:  sr.Name,
						TrackName:  track.Name,
						Result:     true,
					})
				}
			}
		}
	} else {
		// search the album by id was failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the album by ID failed...\t:\t[%s]", id),
		})
	}

	return likeResults
}

// LikeTrackById : returns error if liking is failed
func LikeTrackById(client *spotify.Client, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(client, id); sr.Result && sr.Type == "Track" {
		// execute like track
		if err := client.AddTracksToLibrary(context.Background(), spotify.ID(id)); err != nil {
			likeResults = append(likeResults, &LikeResult{
				// like the track searched by id was failed
				ID:         sr.ID,
				Type:       "Track",
				ArtistName: sr.Artist,
				AlbumName:  sr.Album,
				TrackName:  sr.Name,
				Result:     false,
				Error:      err,
			})
		} else {
			// like the track searched by id was succeeded
			likeResults = append(likeResults, &LikeResult{
				ID:         sr.ID,
				Type:       "Track",
				ArtistName: sr.Artist,
				AlbumName:  sr.Album,
				TrackName:  sr.Name,
				Result:     true,
			})
		}
	} else {
		// search the track by id was failed
		likeResults = append(likeResults, &LikeResult{
			Result:       false,
			Error:        sr.Error,
			ErrorMessage: fmt.Sprintf("Search the track by ID failed...\t:\t[%s]", id),
		})

	}

	return likeResults
}
