package api

import (
	"context"
	"sort"
	"strings"

	"github.com/yanosea/spotlike/util"

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
	// Skip is whether execute like was skipped or not (true : skipped, false : not skipped)
	Skip bool
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

	// skip like
	if !force && alreadyLiked[0] {
		likeResults = append(likeResults, &LikeResult{
			ID:          searchResult.Id,
			Type:        searchResult.Type,
			ArtistNames: searchResult.ArtistNames,
			Skip:        true,
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

		// skip like
		if !force && alreadyLiked[0] {
			likeResults = append(likeResults, &LikeResult{
				ID:          album.ID.String(),
				Type:        spotify.SearchTypeAlbum,
				ArtistNames: combineArtistNames(album.Artists),
				AlbumName:   album.Name,
				Skip:        true,
			})

			continue
		}

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

		// skip like
		if !force && alreadyLiked[0] {
			likeResults = append(likeResults, &LikeResult{
				ID:          track.Track.ID.String(),
				Type:        spotify.SearchTypeTrack,
				ArtistNames: combineArtistNames(track.Track.Artists),
				AlbumName:   track.AlbumName,
				TrackName:   track.Track.Name,
				Skip:        true,
			})

			continue
		}

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
	if err != nil || searchResult.Type != spotify.SearchTypeArtist {
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

	// skip like
	if !force && alreadyLiked[0] {
		likeResults = append(likeResults, &LikeResult{
			ID:          searchResult.Id,
			Type:        searchResult.Type,
			ArtistNames: searchResult.ArtistNames,
			AlbumName:   searchResult.AlbumName,
			Skip:        true,
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
	if err != nil || searchResult.Type != spotify.SearchTypeArtist {
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

		// skip like
		if !force && alreadyLiked[0] {
			likeResults = append(likeResults, &LikeResult{
				ID:          track.ID.String(),
				Type:        spotify.SearchTypeTrack,
				ArtistNames: combineArtistNames(track.Artists),
				AlbumName:   track.Album.Name,
				TrackName:   track.Name,
				Skip:        true,
			})

			continue
		}

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
	if err != nil || searchResult.Type != spotify.SearchTypeArtist {
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

	// skip like
	if !force && alreadyLiked[0] {
		likeResults = append(likeResults, &LikeResult{
			ID:          searchResult.Id,
			Type:        searchResult.Type,
			ArtistNames: searchResult.ArtistNames,
			AlbumName:   searchResult.AlbumName,
			TrackName:   searchResult.TrackName,
			Skip:        true,
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
