/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package api

import (
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// LikeResult : like result from spotify api
type LikeResult struct {
	Type   string
	Name   string
	Result bool
	Error  error
}

// likelikeResults : like results from spotify api
var likeResults []*LikeResult

// LikeArtistById : returns like result
func LikeArtistById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	sr := SearchById(spt, id)
	// execute like artist
	if err := spt.Client.FollowUser(spt.Context, spotify.ID(id)); err != nil {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Artist",
			Name:   sr.Name,
			Result: false,
			Error:  err,
		})
	} else {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Artist",
			Name:   sr.Name,
			Result: false,
			Error:  nil,
		})
	}

	return likeResults
}

// LikeAllAlbumsReleasedByArtistById : returns like results
func LikeAllAlbumsReleasedByArtistById(spt *SpotifyClient, id string) []*LikeResult {
	// get all albums released by the artist
	if allAlbums, err := spt.Client.GetArtistAlbums(spt.Context, spotify.ID(id), nil); err != nil {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Album",
			Name:   "",
			Result: false,
			Error:  err,
		})
	} else {
		// like by album
		for _, album := range allAlbums.Albums {
			if spt.Client.AddAlbumsToLibrary(spt.Context, album.ID); err != nil {
				likeResults = append(likeResults, &LikeResult{
					Type:   "Album",
					Name:   album.Name,
					Result: false,
					Error:  err,
				})
			} else {
				likeResults = append(likeResults, &LikeResult{
					Type:   "Album",
					Name:   album.Name,
					Result: true,
					Error:  nil,
				})
			}
		}
	}

	return likeResults
}

// LikeAllTracksReleasedByArtistById : returns results
func LikeAllTracksReleasedByArtistById(spt *SpotifyClient, id string) []*LikeResult {
	// get all albums released by the artist
	if allAlbums, err := spt.Client.GetArtistAlbums(spt.Context, spotify.ID(id), nil); err != nil {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Track",
			Name:   "",
			Result: false,
			Error:  err,
		})
	} else {
		// get all tracks from all albums
		var allTracks []spotify.SimpleTrack
		for _, album := range allAlbums.Albums {
			if tracks, err := spt.Client.GetAlbumTracks(spt.Context, album.ID); err != nil {
				likeResults = append(likeResults, &LikeResult{
					Type:   "Track",
					Name:   "",
					Result: false,
					Error:  err,
				})
			} else {
				for _, track := range tracks.Tracks {
					allTracks = append(allTracks, track)
				}
			}
		}

		// like by track
		for _, track := range allTracks {
			if spt.Client.AddTracksToLibrary(spt.Context, track.ID); err != nil {
				likeResults = append(likeResults, &LikeResult{
					Type:   "Track",
					Name:   track.Name,
					Result: false,
					Error:  err,
				})
			} else {
				likeResults = append(likeResults, &LikeResult{
					Type:   "Track",
					Name:   track.Name,
					Result: true,
					Error:  nil,
				})
			}
		}
	}

	return likeResults
}

// LikeAlbumById : returns error if liking is failed
func LikeAlbumById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	sr := SearchById(spt, id)
	// execute like album
	if err := spt.Client.AddAlbumsToLibrary(spt.Context, spotify.ID(id)); err != nil {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Album",
			Name:   sr.Name,
			Result: false,
			Error:  err,
		})
	} else {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Album",
			Name:   sr.Name,
			Result: true,
			Error:  nil,
		})
	}

	return likeResults
}

// LikeAllTracksOnAlbumById : returns error if liking is failed
func LikeAllTracksOnAlbumById(spt *SpotifyClient, id string) []*LikeResult {
	// get all tracks on the artist
	if allTracks, err := spt.Client.GetAlbumTracks(spt.Context, spotify.ID(id)); err != nil {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Track",
			Name:   "",
			Result: false,
			Error:  err,
		})
	} else {
		for _, track := range allTracks.Tracks {
			// execute like track
			if spt.Client.AddTracksToLibrary(spt.Context, track.ID); err != nil {
				likeResults = append(likeResults, &LikeResult{
					Type:   "Track",
					Name:   track.Name,
					Result: false,
					Error:  err,
				})
			} else {
				likeResults = append(likeResults, &LikeResult{
					Type:   "Track",
					Name:   track.Name,
					Result: true,
					Error:  nil,
				})
			}
		}
	}

	return likeResults
}

// LikeTrackById : returns error if liking is failed
func LikeTrackById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	sr := SearchById(spt, id)
	// execute like track
	if err := spt.Client.AddTracksToLibrary(spt.Context, spotify.ID(id)); err != nil {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Track",
			Name:   sr.Name,
			Result: false,
			Error:  err,
		})
	} else {
		likeResults = append(likeResults, &LikeResult{
			Type:   "Track",
			Name:   sr.Name,
			Result: true,
			Error:  nil,
		})
	}

	return likeResults
}
