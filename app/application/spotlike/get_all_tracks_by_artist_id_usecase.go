package spotlike

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
)

// getAllTracksByArtistIdUseCase is a struct that contains the use case of getting for an artist.
type getAllTracksByArtistIdUseCase struct {
	trackRepo trackDomain.TrackRepository
}

// NewGetAllTracksByArtistIdUseCase returns a new instance of the GetAllTracksByArtistIdUseCase struct.
func NewGetAllTracksByArtistIdUseCase(tracksRepo trackDomain.TrackRepository) *getAllTracksByArtistIdUseCase {
	return &getAllTracksByArtistIdUseCase{
		trackRepo: tracksRepo,
	}
}

// GetAllTracksByArtistIdUseCaseOutputDto is a DTO struct that contains the output data of the getAllTracksUseCase.
type GetAllTracksByArtistIdUseCaseOutputDto struct {
	ID          string
	Artists     string
	Album       string
	Name        string
	TrackNumber spotify.Numeric
	ReleaseDate time.Time
}

// Run returns the get result of the tracks.
func (uc *getAllTracksByArtistIdUseCase) Run(ctx context.Context, id string) ([]*GetAllTracksByArtistIdUseCaseOutputDto, error) {
	tracks, err := uc.trackRepo.FindByArtistId(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}

	// group tracks by album
	albumMap := make(map[string][]*trackDomain.Track)
	for _, track := range tracks {
		albumID := track.Album.ID.String()
		albumMap[albumID] = append(albumMap[albumID], track)
	}

	// sort albums by release date
	var albumIDs []string
	for albumID := range albumMap {
		albumIDs = append(albumIDs, albumID)
	}
	sort.Slice(albumIDs, func(i, j int) bool {
		return albumMap[albumIDs[i]][0].ReleaseDate.Before(albumMap[albumIDs[j]][0].ReleaseDate)
	})

	var getAllTracksByArtistIdUseCaseOutputDtos []*GetAllTracksByArtistIdUseCaseOutputDto
	for _, albumID := range albumIDs {
		albumTracks := albumMap[albumID]
		// sort tracks by track number
		sort.Slice(albumTracks, func(i, j int) bool {
			return albumTracks[i].TrackNumber < albumTracks[j].TrackNumber
		})

		for _, track := range albumTracks {
			artistNames := make([]string, len(track.Artists))
			for i, artist := range track.Artists {
				artistNames[i] = artist.Name
			}

			getAllTracksByArtistIdUseCaseOutputDto := &GetAllTracksByArtistIdUseCaseOutputDto{
				ID:          track.ID.String(),
				Artists:     strings.Join(artistNames, ", "),
				Album:       track.Album.Name,
				Name:        track.Name,
				TrackNumber: track.TrackNumber,
				ReleaseDate: track.ReleaseDate,
			}
			getAllTracksByArtistIdUseCaseOutputDtos = append(getAllTracksByArtistIdUseCaseOutputDtos, getAllTracksByArtistIdUseCaseOutputDto)
		}
	}

	return getAllTracksByArtistIdUseCaseOutputDtos, nil
}
