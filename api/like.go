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
	// ID : content's ID
	ID string
	// Type : content type
	Type string
	// Name : content name
	Name string
	// Result : like result (true : succeeded / false : failed)
	Result bool
	// Error : error returned from spotify api
	Error error
	// Message : message for like result
	Message string
}

// likelikeResults : like results from spotify api
var likeResults []*LikeResult

// LikeArtistById : returns like result
func LikeArtistById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(spt, id); sr.Result && sr.Type == "Artist" {
		// execute like artist
		if err := spt.Client.FollowUser(spt.Context, spotify.ID(sr.ID)); err != nil {
			likeResults = append(likeResults, &LikeResult{
				ID:      sr.ID,
				Type:    "Artist",
				Name:    sr.Name,
				Result:  false,
				Error:   err,
				Message: "Like the artist failed...",
			})
		} else {
			likeResults = append(likeResults, &LikeResult{
				ID:      sr.ID,
				Type:    "Artist",
				Name:    sr.Name,
				Result:  true,
				Error:   nil,
				Message: "Like the artist succeeded!",
			})
		}
	} else {
		likeResults = append(likeResults, &LikeResult{
			ID:      id,
			Type:    "",
			Name:    "",
			Result:  false,
			Error:   sr.Error,
			Message: "Search the artist by ID failed...",
		})
	}

	return likeResults
}

// LikeAllAlbumsReleasedByArtistById : returns like results
func LikeAllAlbumsReleasedByArtistById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(spt, id); sr.Result && sr.Type == "Artist" {
		// get all albums released by the artist
		if allAlbums, err := spt.Client.GetArtistAlbums(spt.Context, spotify.ID(sr.ID), nil); err != nil {
			likeResults = append(likeResults, &LikeResult{
				ID:      sr.ID,
				Type:    "",
				Name:    "",
				Result:  false,
				Error:   err,
				Message: "Get all albums by the artist searched by ID failed...",
			})
		} else {
			// like by album
			for _, album := range allAlbums.Albums {
				if spt.Client.AddAlbumsToLibrary(spt.Context, album.ID); err != nil {
					likeResults = append(likeResults, &LikeResult{
						ID:      album.ID.String(),
						Type:    "Album",
						Name:    album.Name,
						Result:  false,
						Error:   err,
						Message: "Like the album by the artist searched by ID failed...",
					})
				} else {
					likeResults = append(likeResults, &LikeResult{
						ID:      album.ID.String(),
						Type:    "Album",
						Name:    album.Name,
						Result:  true,
						Error:   nil,
						Message: "Like the album succeeded!",
					})
				}
			}
		}
	} else {
		likeResults = append(likeResults, &LikeResult{
			ID:      id,
			Type:    "",
			Name:    "",
			Result:  false,
			Error:   sr.Error,
			Message: "Search the artist by ID failed...",
		})
	}

	return likeResults
}

// LikeAllTracksReleasedByArtistById : returns results
func LikeAllTracksReleasedByArtistById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(spt, id); sr.Result && sr.Type == "Artist" {
		// get all albums released by the artist
		if allAlbums, err := spt.Client.GetArtistAlbums(spt.Context, spotify.ID(sr.ID), nil); err != nil {
			likeResults = append(likeResults, &LikeResult{
				ID:      sr.ID,
				Type:    "",
				Name:    "",
				Result:  false,
				Error:   err,
				Message: "Get all albums by the artist searched by ID failed...",
			})
		} else {
			// get all tracks from all albums
			var allTracks []spotify.SimpleTrack
			for _, album := range allAlbums.Albums {
				if tracks, err := spt.Client.GetAlbumTracks(spt.Context, album.ID); err != nil {
					likeResults = append(likeResults, &LikeResult{
						ID:      id,
						Type:    "",
						Name:    "",
						Result:  false,
						Error:   err,
						Message: "Get all tracks in all albums by the artist searched by ID failed...",
					})
				} else {
					for _, track := range tracks.Tracks {
						allTracks = append(allTracks, track)
					}
					// like by track
					for _, track := range allTracks {
						if spt.Client.AddTracksToLibrary(spt.Context, track.ID); err != nil {
							likeResults = append(likeResults, &LikeResult{
								ID:      track.ID.String(),
								Type:    "Track",
								Name:    track.Name,
								Result:  false,
								Error:   err,
								Message: "Like the track in the album by the artist searched by ID failed...",
							})
						} else {
							likeResults = append(likeResults, &LikeResult{
								ID:      track.ID.String(),
								Type:    "Track",
								Name:    track.Name,
								Result:  true,
								Error:   nil,
								Message: "Like the track succeeded!",
							})
						}
					}
				}
			}
		}
	} else {
		likeResults = append(likeResults, &LikeResult{
			ID:      id,
			Type:    "",
			Name:    "",
			Result:  false,
			Error:   sr.Error,
			Message: "Search the artist by ID failed...",
		})
	}

	return likeResults
}

// LikeAlbumById : returns error if liking is failed
func LikeAlbumById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(spt, id); sr.Result && sr.Type == "Album" {
		// execute like album
		if err := spt.Client.AddAlbumsToLibrary(spt.Context, spotify.ID(id)); err != nil {
			likeResults = append(likeResults, &LikeResult{
				ID:      id,
				Type:    "Album",
				Name:    sr.Name,
				Result:  false,
				Error:   err,
				Message: "Like the album failed...",
			})
		} else {
			likeResults = append(likeResults, &LikeResult{
				ID:      id,
				Type:    "Album",
				Name:    sr.Name,
				Result:  true,
				Error:   nil,
				Message: "Like the album succeeded!",
			})
		}
	} else {
		likeResults = append(likeResults, &LikeResult{
			ID:      id,
			Type:    "",
			Name:    "",
			Result:  false,
			Error:   sr.Error,
			Message: "Search the album by ID failed...",
		})
	}

	return likeResults
}

// LikeAllTracksOnAlbumById : returns error if liking is failed
func LikeAllTracksOnAlbumById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(spt, id); sr.Result && sr.Type == "ALbum" {
		// get all tracks on the artist
		if allTracks, err := spt.Client.GetAlbumTracks(spt.Context, spotify.ID(id)); err != nil {
			likeResults = append(likeResults, &LikeResult{
				ID:      sr.ID,
				Type:    "",
				Name:    "",
				Result:  false,
				Error:   err,
				Message: "Get all tracks in the album searched by ID failed...",
			})
		} else {
			for _, track := range allTracks.Tracks {
				// execute like track
				if spt.Client.AddTracksToLibrary(spt.Context, track.ID); err != nil {
					likeResults = append(likeResults, &LikeResult{
						ID:      track.ID.String(),
						Type:    "Track",
						Name:    track.Name,
						Result:  false,
						Error:   err,
						Message: "Like the track in the album searched by ID failed...",
					})
				} else {
					likeResults = append(likeResults, &LikeResult{
						ID:      track.ID.String(),
						Type:    "Track",
						Name:    track.Name,
						Result:  true,
						Error:   nil,
						Message: "Like the track succeeded!",
					})
				}
			}
		}
	} else {
		likeResults = append(likeResults, &LikeResult{
			ID:      id,
			Type:    "",
			Name:    "",
			Result:  false,
			Error:   sr.Error,
			Message: "Search the album by ID failed...",
		})
	}

	return likeResults
}

// LikeTrackById : returns error if liking is failed
func LikeTrackById(spt *SpotifyClient, id string) []*LikeResult {
	// execute search by id
	if sr := SearchById(spt, id); sr.Result && sr.Type == "Track" {
		// execute like track
		if err := spt.Client.AddTracksToLibrary(spt.Context, spotify.ID(id)); err != nil {
			likeResults = append(likeResults, &LikeResult{
				ID:      sr.ID,
				Type:    "Track",
				Name:    sr.Name,
				Result:  false,
				Error:   err,
				Message: "Like the track failed...",
			})
		} else {
			likeResults = append(likeResults, &LikeResult{
				ID:      sr.ID,
				Type:    "Track",
				Name:    sr.Name,
				Result:  true,
				Error:   nil,
				Message: "Like the track succeeded!",
			})
		}
	} else {
		likeResults = append(likeResults, &LikeResult{
			ID:      id,
			Type:    "",
			Name:    "",
			Result:  false,
			Error:   sr.Error,
			Message: "Search the track by ID failed...",
		})

	}

	return likeResults
}
