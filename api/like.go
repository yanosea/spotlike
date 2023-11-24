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
	// Result is the like result (true : succeeded, false : failed)
	Result bool
	// Skip is whether execute like was skipped or not (true : skipped, false : not skipped)
	Skip bool
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
func LikeArtistById(client *spotify.Client, id string, force bool) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Artist" {
		// check the artist has been already liked
		if alreadyLiked, err := client.CurrentUserFollows(context.Background(), "artist", spotify.ID(sr.ID)); err == nil {
			if !force && alreadyLiked[0] {
				// like skipped
				likeResults = append(likeResults, &LikeResult{
					ID:          sr.ID,
					Type:        "Artist",
					ArtistNames: sr.ArtistNames,
					Result:      true,
					Skip:        true,
				})
			} else if err := client.FollowArtist(context.Background(), spotify.ID(sr.ID)); err != nil {
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
			// check failed
			likeResults = append(likeResults, &LikeResult{
				ID:           sr.ID,
				Type:         "Artist",
				ArtistNames:  sr.ArtistNames,
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Check whether the artist has been already liked failed...\t:\t[%s]", sr.ID),
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
func LikeAllAlbumsReleasedByArtistById(client *spotify.Client, id string, force bool) []*LikeResult {
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

			for _, album := range allAlbums.Albums {
				// check the album has been already liked
				if alreadyLiked, err := client.UserHasAlbums(context.Background(), album.ID); err == nil {
					if !force && alreadyLiked[0] {
						// like skipped
						likeResults = append(likeResults, &LikeResult{
							ID:          album.ID.String(),
							Type:        "Album",
							ArtistNames: combineArtistNames(album.Artists),
							AlbumName:   album.Name,
							Result:      true,
							Skip:        true,
						})
					} else if err := client.AddAlbumsToLibrary(context.Background(), album.ID); err != nil {
						likeResults = append(likeResults, &LikeResult{
							// like failed
							ID:          album.ID.String(),
							Type:        "Album",
							ArtistNames: combineArtistNames(album.Artists),
							AlbumName:   album.Name,
							Result:      false,
							Error:       err,
						})
					} else {
						// like succeeded
						likeResults = append(likeResults, &LikeResult{
							ID:          album.ID.String(),
							Type:        "Album",
							ArtistNames: combineArtistNames(album.Artists),
							AlbumName:   album.Name,
							Result:      true,
						})
					}
				} else {
					// check failed
					likeResults = append(likeResults, &LikeResult{
						ID:           album.ID.String(),
						Type:         "Artist",
						ArtistNames:  combineArtistNames(album.Artists),
						AlbumName:    album.Name,
						Result:       false,
						Error:        err,
						ErrorMessage: fmt.Sprintf("Check whether the album has been already liked failed...\t:\t[%s]", album.ID.String()),
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
func LikeAllTracksReleasedByArtistById(client *spotify.Client, id string, force bool) []*LikeResult {
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

			for _, track := range allTracks {
				// check the track has been already liked
				if alreadyLiked, err := client.UserHasTracks(context.Background(), track.Track.ID); err == nil {
					if !force && alreadyLiked[0] {
						// like skipped
						likeResults = append(likeResults, &LikeResult{
							ID:          track.Track.ID.String(),
							Type:        "Track",
							ArtistNames: combineArtistNames(track.Track.Artists),
							AlbumName:   track.AlbumName,
							TrackName:   track.Track.Name,
							Result:      true,
							Skip:        true,
						})
					} else if err := client.AddTracksToLibrary(context.Background(), track.Track.ID); err != nil {
						// like failed
						likeResults = append(likeResults, &LikeResult{
							ID:          track.Track.ID.String(),
							Type:        "Track",
							ArtistNames: combineArtistNames(track.Track.Artists),
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
							ArtistNames: combineArtistNames(track.Track.Artists),
							AlbumName:   track.AlbumName,
							TrackName:   track.Track.Name,
							Result:      true,
						})
					}
				} else {
					// check failed
					likeResults = append(likeResults, &LikeResult{
						ID:           track.Track.ID.String(),
						Type:         "Track",
						ArtistNames:  combineArtistNames(track.Track.Artists),
						AlbumName:    track.AlbumName,
						TrackName:    track.Track.Name,
						Result:       false,
						Error:        err,
						ErrorMessage: fmt.Sprintf("Check whether the track has been already liked failed...\t:\t[%s]", track.Track.ID.String()),
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
func LikeAlbumById(client *spotify.Client, id string, force bool) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Album" {
		// check the album has been already liked
		if alreadyLiked, err := client.UserHasAlbums(context.Background(), spotify.ID(sr.ID)); err == nil {
			if !force && alreadyLiked[0] {
				// like skipped
				likeResults = append(likeResults, &LikeResult{
					ID:          sr.ID,
					Type:        "Album",
					ArtistNames: sr.ArtistNames,
					AlbumName:   sr.AlbumName,
					Result:      true,
					Skip:        true,
				})
			} else if err := client.AddAlbumsToLibrary(context.Background(), spotify.ID(sr.ID)); err != nil {
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
			// check failed
			likeResults = append(likeResults, &LikeResult{
				ID:           sr.ID,
				Type:         "Album",
				ArtistNames:  sr.ArtistNames,
				AlbumName:    sr.AlbumName,
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Check whether the album has been already liked failed...\t:\t[%s]", sr.ID),
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
func LikeAllTracksInAlbumById(client *spotify.Client, id string, force bool) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Album" {
		// get all tracks in the album
		if allTracks, err := client.GetAlbumTracks(context.Background(), spotify.ID(sr.ID)); err != nil {
			// getting all tracks in the album searched by ID failed
			likeResults = append(likeResults, &LikeResult{
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Get all tracks in the album searched by ID failed...\t:\t[%s]", id),
			})
		} else {
			for _, track := range allTracks.Tracks {
				// check the track has been already liked
				if alreadyLiked, err := client.UserHasTracks(context.Background(), track.ID); err == nil {
					if !force && alreadyLiked[0] {
						// like skipped
						likeResults = append(likeResults, &LikeResult{
							ID:          track.ID.String(),
							Type:        "Track",
							ArtistNames: combineArtistNames(track.Artists),
							AlbumName:   track.Album.Name,
							TrackName:   track.Name,
							Result:      true,
							Skip:        true,
						})
					} else if err := client.AddTracksToLibrary(context.Background(), track.ID); err != nil {
						// like failed
						likeResults = append(likeResults, &LikeResult{
							ID:          track.ID.String(),
							Type:        "Track",
							ArtistNames: combineArtistNames(track.Artists),
							AlbumName:   track.Album.Name,
							TrackName:   track.Name,
							Result:      false,
							Error:       err,
						})
					} else {
						// like succeeded
						likeResults = append(likeResults, &LikeResult{
							ID:          track.ID.String(),
							Type:        "Track",
							ArtistNames: combineArtistNames(track.Artists),
							AlbumName:   track.Album.Name,
							TrackName:   track.Name,
							Result:      true,
						})
					}
				} else {
					// check failed
					likeResults = append(likeResults, &LikeResult{
						ID:           track.ID.String(),
						Type:         "Track",
						ArtistNames:  combineArtistNames(track.Artists),
						AlbumName:    track.Album.Name,
						TrackName:    track.Name,
						Result:       false,
						Error:        err,
						ErrorMessage: fmt.Sprintf("Check whether the track has been already liked failed...\t:\t[%s]", track.ID.String()),
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
func LikeTrackById(client *spotify.Client, id string, force bool) []*LikeResult {
	// execute search
	if sr := SearchById(client, id); sr.Result && sr.Type == "Track" {
		// check the track has been already liked
		if alreadyLiked, err := client.UserHasTracks(context.Background(), spotify.ID(sr.ID)); err == nil {
			if !force && alreadyLiked[0] {
				// like skipped
				likeResults = append(likeResults, &LikeResult{
					ID:          sr.ID,
					Type:        "Track",
					ArtistNames: sr.ArtistNames,
					AlbumName:   sr.AlbumName,
					TrackName:   sr.TrackName,
					Result:      true,
					Skip:        true,
				})
			} else if err := client.AddTracksToLibrary(context.Background(), spotify.ID(sr.ID)); err != nil {
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
			// check failed
			likeResults = append(likeResults, &LikeResult{
				ID:           sr.ID,
				Type:         "Track",
				ArtistNames:  sr.ArtistNames,
				AlbumName:    sr.AlbumName,
				TrackName:    sr.TrackName,
				Result:       false,
				Error:        err,
				ErrorMessage: fmt.Sprintf("Check whether the track has been already liked failed...\t:\t[%s]", sr.ID),
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
