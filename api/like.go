package api

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/yanosea/spotlike/util"

	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// LikeResult represents the result of liking content using the Spotify API.
type LikeResult struct {
	// ID is content's id
	ID string
	// Type is the content type
	Type spotify.SearchType
	// ArtistNames is the names of the artists
	ArtistNames string
	// AlbumName is the album name
	AlbumName string
	// TrackName is the track name
	TrackName string
	// AlreadyLiked is whether already liked or not (true : already liked, false : not already liked)
	AlreadyLiked bool
	// Refused is whether answer "y" or not (true : refused, false : not refused)
	Refused bool
	// Error is the error returned from the Spotify API
	Error error
}

// TrackWithAlbumName represents a Spotify simple track with an album name.
type TrackWithAlbumName struct {
	// Track is Spotify simple track
	Track spotify.SimpleTrack
	// AlbumName is the Name of the album the track is included in
	AlbumName string
}

const (
	// input_label_confirm_like_artist is the message confirming like the artist.
	input_label_confirm_like_artist = `Do you execute like artist [%s]]`
	// input_label_confirm_like_artist is the message confirming like the album.
	input_label_confirm_like_album = `Do you execute like album "[%s]" by "[%s]"`
	// input_label_confirm_like_artist is the messagr confirming like the track.
	input_label_confirm_like_track = `Do you execute like track "[%s]" in "[%s]" by "[%s]"`
)

// LikeArtistById returns the like result for an artist with the given ID.
func LikeArtistById(client *spotify.Client, id string, force bool) []*LikeResult {
	var likeResults []*LikeResult
	// execute search
	searchResult, err := SearchById(client, id)
	if err != nil || searchResult.Type != spotify.SearchTypeArtist {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// check the artist has been already liked
	alreadyLiked, err := client.CurrentUserFollows(context.Background(), strings.ToLower(util.STRING_ARTIST), spotify.ID(searchResult.Id))
	if err != nil {
		// check failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// skip like if already liked
	if alreadyLiked[0] {
		likeResults = append(likeResults, &LikeResult{
			ID:           searchResult.Id,
			Type:         searchResult.Type,
			ArtistNames:  searchResult.ArtistNames,
			AlreadyLiked: true,
			Refused:      false,
		})

		return likeResults
	}
	// confirm like
	answer := "y"
	if !force {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf(input_label_confirm_like_artist, searchResult.ArtistNames),
			IsConfirm: true,
		}

		input, err := prompt.Run()
		if err != nil {
			answer = "n"
		}

		answer = input
	}
	// skip like if refused
	if answer == "n" {
		likeResults = append(likeResults, &LikeResult{
			ID:           searchResult.Id,
			Type:         searchResult.Type,
			ArtistNames:  searchResult.ArtistNames,
			AlreadyLiked: false,
			Refused:      true,
		})

		return likeResults
	}
	// execute like
	if err := client.FollowArtist(context.Background(), spotify.ID(searchResult.Id)); err != nil {
		// like failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// like succeeded
	likeResults = append(likeResults, &LikeResult{
		ID:          searchResult.Id,
		Type:        searchResult.Type,
		ArtistNames: searchResult.ArtistNames,
	})

	return likeResults
}

// LikeAllAlbumsReleasedByArtistById returns the like results for all albums released by an artist with the given ID.
func LikeAllAlbumsReleasedByArtistById(client *spotify.Client, id string, force bool) []*LikeResult {
	var likeResults []*LikeResult
	// execute search
	searchResult, err := SearchById(client, id)
	if err != nil || searchResult.Type != spotify.SearchTypeArtist {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// get all albums released by the artist
	allAlbums, err := client.GetArtistAlbums(context.Background(), spotify.ID(searchResult.Id), nil)
	if err != nil {
		likeResults = append(likeResults, &LikeResult{
			// getting all albums by the artist searched by ID failed
			Error: err,
		})
		return likeResults
	}
	// sort albums by release date
	sort.Slice(allAlbums.Albums, func(i, j int) bool {
		return allAlbums.Albums[i].ReleaseDateTime().Before(allAlbums.Albums[j].ReleaseDateTime())
	})
	// for each album
	for _, album := range allAlbums.Albums {
		// check the album has been already liked
		alreadyLiked, err := client.UserHasAlbums(context.Background(), album.ID)
		if err != nil {
			// check failed
			likeResults = append(likeResults, &LikeResult{
				Error: err,
			})

			continue
		}
		// skip like if already liked
		if alreadyLiked[0] {
			likeResults = append(likeResults, &LikeResult{
				ID:           album.ID.String(),
				Type:         spotify.SearchTypeAlbum,
				ArtistNames:  combineArtistNames(album.Artists),
				AlbumName:    album.Name,
				AlreadyLiked: true,
				Refused:      false,
			})

			continue
		}
		// confirm like
		answer := "y"
		if !force {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf(input_label_confirm_like_album, album.Name, combineArtistNames(album.Artists)),
				IsConfirm: true,
			}

			input, err := prompt.Run()
			if err != nil {
				answer = "n"
			}

			answer = input
		}
		// skip like if refused
		if answer == "n" {
			likeResults = append(likeResults, &LikeResult{
				ID:           album.ID.String(),
				Type:         spotify.SearchTypeAlbum,
				ArtistNames:  combineArtistNames(album.Artists),
				AlbumName:    album.Name,
				AlreadyLiked: false,
				Refused:      true,
			})

			continue
		}
		// execute like
		if err := client.AddAlbumsToLibrary(context.Background(), album.ID); err != nil {
			// like failed
			likeResults = append(likeResults, &LikeResult{
				Error: err,
			})

			continue
		}
		// like succeeded
		likeResults = append(likeResults, &LikeResult{
			ID:          album.ID.String(),
			Type:        spotify.SearchTypeAlbum,
			ArtistNames: combineArtistNames(album.Artists),
			AlbumName:   album.Name,
		})
	}

	return likeResults
}

// LikeAllTracksReleasedByArtistById returns the like results for all tracks released by an artist with the given ID.
func LikeAllTracksReleasedByArtistById(client *spotify.Client, id string, force bool) []*LikeResult {
	var likeResults []*LikeResult
	// execute search
	searchResult, err := SearchById(client, id)
	if err != nil || searchResult.Type != spotify.SearchTypeArtist {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// get all albums released by the artist
	allAlbums, err := client.GetArtistAlbums(context.Background(), spotify.ID(searchResult.Id), nil)
	if err != nil {
		likeResults = append(likeResults, &LikeResult{
			// getting all albums by the artist searched by ID failed
			Error: err,
		})
		return likeResults
	}
	// sort albums by release date
	sort.Slice(allAlbums.Albums, func(i, j int) bool {
		return allAlbums.Albums[i].ReleaseDateTime().Before(allAlbums.Albums[j].ReleaseDateTime())
	})
	// get all tracks from all albums
	var allTracks []TrackWithAlbumName
	for _, album := range allAlbums.Albums {
		tracks, err := client.GetAlbumTracks(context.Background(), album.ID)
		if err != nil {
			// getting all tracks in all albums by the artist searched by ID failed
			likeResults = append(likeResults, &LikeResult{
				Error: err,
			})
		}
		// add album name
		for _, track := range tracks.Tracks {
			trackWithAlbumName := &TrackWithAlbumName{
				Track:     track,
				AlbumName: album.Name,
			}

			allTracks = append(allTracks, *trackWithAlbumName)
		}
	}
	// for each track
	for _, track := range allTracks {
		// check the track has been already liked
		alreadyLiked, err := client.UserHasTracks(context.Background(), track.Track.ID)
		if err != nil {
			// check failed
			likeResults = append(likeResults, &LikeResult{
				Error: err,
			})

			continue
		}
		// skip like if already liked
		if alreadyLiked[0] {
			likeResults = append(likeResults, &LikeResult{
				ID:           track.Track.ID.String(),
				Type:         spotify.SearchTypeTrack,
				ArtistNames:  combineArtistNames(track.Track.Artists),
				AlbumName:    track.AlbumName,
				TrackName:    track.Track.Name,
				AlreadyLiked: true,
				Refused:      false,
			})

			continue
		}
		// confirm like
		answer := "y"
		if !force {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf(input_label_confirm_like_track, track.Track.Name, track.AlbumName, combineArtistNames(track.Track.Artists)),
				IsConfirm: true,
			}

			input, err := prompt.Run()
			if err != nil {
				answer = "n"
			}

			answer = input
		}
		// skip like if refused
		if answer == "n" {
			likeResults = append(likeResults, &LikeResult{
				ID:           track.Track.ID.String(),
				Type:         spotify.SearchTypeTrack,
				ArtistNames:  combineArtistNames(track.Track.Artists),
				AlbumName:    track.AlbumName,
				TrackName:    track.Track.Name,
				AlreadyLiked: false,
				Refused:      true,
			})

			continue
		}
		// execute like
		if err := client.AddTracksToLibrary(context.Background(), track.Track.ID); err != nil {
			// like failed
			likeResults = append(likeResults, &LikeResult{
				Error: err,
			})

			continue
		}
		// like succeeded
		likeResults = append(likeResults, &LikeResult{
			ID:          track.Track.ID.String(),
			Type:        spotify.SearchTypeTrack,
			ArtistNames: combineArtistNames(track.Track.Artists),
			AlbumName:   track.AlbumName,
			TrackName:   track.Track.Name,
		})
	}

	return likeResults
}

// LikeAlbumById returns an error if liking an album with the given ID is failed.
func LikeAlbumById(client *spotify.Client, id string, force bool) []*LikeResult {
	var likeResults []*LikeResult
	// execute search
	searchResult, err := SearchById(client, id)
	if err != nil || searchResult.Type != spotify.SearchTypeAlbum {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// check the album has been already liked
	alreadyLiked, err := client.UserHasAlbums(context.Background(), spotify.ID(searchResult.Id))
	if err != nil {
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// skip like if already liked
	if alreadyLiked[0] {
		likeResults = append(likeResults, &LikeResult{
			ID:           searchResult.Id,
			Type:         searchResult.Type,
			ArtistNames:  searchResult.ArtistNames,
			AlbumName:    searchResult.AlbumName,
			AlreadyLiked: true,
			Refused:      false,
		})

		return likeResults
	}
	// confirm like
	answer := "y"
	if !force {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf(input_label_confirm_like_album, searchResult.AlbumName, searchResult.ArtistNames),
			IsConfirm: true,
		}

		input, err := prompt.Run()
		if err != nil {
			answer = "n"
		}

		answer = input
	}
	// skip like if refused
	if answer == "n" {
		likeResults = append(likeResults, &LikeResult{
			ID:           searchResult.Id,
			Type:         searchResult.Type,
			ArtistNames:  searchResult.ArtistNames,
			AlbumName:    searchResult.AlbumName,
			AlreadyLiked: false,
			Refused:      true,
		})

		return likeResults
	}
	// execute like
	if err := client.AddAlbumsToLibrary(context.Background(), spotify.ID(searchResult.Id)); err != nil {
		// like failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})
	}
	// like succeeded
	likeResults = append(likeResults, &LikeResult{
		ID:          id,
		Type:        searchResult.Type,
		ArtistNames: searchResult.ArtistNames,
		AlbumName:   searchResult.AlbumName,
	})

	return likeResults
}

// LikeAllTracksInAlbumById returns an error if liking all tracks in an album with the given ID is failed.
func LikeAllTracksInAlbumById(client *spotify.Client, id string, force bool) []*LikeResult {
	var likeResults []*LikeResult
	// execute search
	searchResult, err := SearchById(client, id)
	if err != nil || searchResult.Type != spotify.SearchTypeAlbum {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// get all tracks from the albums
	allTracks, err := client.GetAlbumTracks(context.Background(), spotify.ID(searchResult.Id))
	if err != nil {
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})
	}
	// for each track
	for _, track := range allTracks.Tracks {
		// check the track has been already liked
		alreadyLiked, err := client.UserHasTracks(context.Background(), track.ID)
		if err != nil {
			// check failed
			likeResults = append(likeResults, &LikeResult{
				Error: err,
			})

			continue
		}
		// skip like if already liked
		if alreadyLiked[0] {
			likeResults = append(likeResults, &LikeResult{
				ID:           track.ID.String(),
				Type:         spotify.SearchTypeTrack,
				ArtistNames:  combineArtistNames(track.Artists),
				AlbumName:    track.Album.Name,
				TrackName:    track.Name,
				AlreadyLiked: true,
				Refused:      false,
			})

			continue
		}
		// confirm like
		answer := "y"
		if !force {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf(input_label_confirm_like_track, track.Name, track.Album.Name, combineArtistNames(track.Artists)),
				IsConfirm: true,
			}

			input, err := prompt.Run()
			if err != nil {
				answer = "n"
			}

			answer = input
		}
		// skip like if refused
		if answer == "n" {
			likeResults = append(likeResults, &LikeResult{
				ID:           track.ID.String(),
				Type:         spotify.SearchTypeTrack,
				ArtistNames:  combineArtistNames(track.Artists),
				AlbumName:    track.Album.Name,
				TrackName:    track.Name,
				AlreadyLiked: false,
				Refused:      true,
			})

			continue
		}
		// execute like
		if err := client.AddTracksToLibrary(context.Background(), track.ID); err != nil {
			// like failed
			likeResults = append(likeResults, &LikeResult{
				Error: err,
			})

			continue
		}
		// like succeeded
		likeResults = append(likeResults, &LikeResult{
			ID:          track.ID.String(),
			Type:        spotify.SearchTypeTrack,
			ArtistNames: combineArtistNames(track.Artists),
			AlbumName:   track.Album.Name,
			TrackName:   track.Name,
		})
	}

	return likeResults
}

// LikeTrackById returns an error if liking a track with the given ID is failed.
func LikeTrackById(client *spotify.Client, id string, force bool) []*LikeResult {
	var likeResults []*LikeResult
	// execute search
	searchResult, err := SearchById(client, id)
	if err != nil || searchResult.Type != spotify.SearchTypeTrack {
		// search failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// check the track has been already liked
	alreadyLiked, err := client.UserHasTracks(context.Background(), spotify.ID(searchResult.Id))
	if err != nil {
		// check failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// skip like if already liked
	if alreadyLiked[0] {
		likeResults = append(likeResults, &LikeResult{
			ID:           searchResult.Id,
			Type:         searchResult.Type,
			ArtistNames:  searchResult.ArtistNames,
			AlbumName:    searchResult.AlbumName,
			TrackName:    searchResult.TrackName,
			AlreadyLiked: true,
			Refused:      false,
		})

		return likeResults
	}
	// confirm like
	answer := "y"
	if !force {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf(input_label_confirm_like_track, searchResult.TrackName, searchResult.AlbumName, searchResult.ArtistNames),
			IsConfirm: true,
		}

		input, err := prompt.Run()
		if err != nil {
			answer = "n"
		}

		answer = input
	}
	// skip like if refused
	if answer == "n" {
		likeResults = append(likeResults, &LikeResult{
			ID:           searchResult.Id,
			Type:         searchResult.Type,
			ArtistNames:  searchResult.ArtistNames,
			AlbumName:    searchResult.AlbumName,
			TrackName:    searchResult.TrackName,
			AlreadyLiked: false,
			Refused:      true,
		})

		return likeResults
	}
	// execute like
	if err := client.AddTracksToLibrary(context.Background(), spotify.ID(searchResult.Id)); err != nil {
		// like failed
		likeResults = append(likeResults, &LikeResult{
			Error: err,
		})

		return likeResults
	}
	// like succeeded
	likeResults = append(likeResults, &LikeResult{
		ID:          searchResult.Id,
		Type:        searchResult.Type,
		ArtistNames: searchResult.ArtistNames,
		AlbumName:   searchResult.AlbumName,
		TrackName:   searchResult.TrackName,
	})

	return likeResults
}
