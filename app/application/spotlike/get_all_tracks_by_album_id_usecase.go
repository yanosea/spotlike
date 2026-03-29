package spotlike

import (
	"context"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
)

// getAllTracksByAlbumIdUseCase is a struct that contains the use case of getting for an track.
type getAllTracksByAlbumIdUseCase struct {
	trackRepo trackDomain.TrackRepository
}

// NewGetAllTracksByAlbumIdUseCase returns a new instance of the GetAllTracksByAlbumIdUseCase struct.
func NewGetAllTracksByAlbumIdUseCase(trackRepo trackDomain.TrackRepository) *getAllTracksByAlbumIdUseCase {
	return &getAllTracksByAlbumIdUseCase{
		trackRepo: trackRepo,
	}
}

// GetAllTracksByAlbumIdUseCaseOutputDto is a DTO struct that contains the output data of the getAllTracksUseCase.
type GetAllTracksByAlbumIdUseCaseOutputDto struct {
	ID          string
	Artists     string
	Album       string
	Name        string
	TrackNumber spotify.Numeric
	ReleaseDate time.Time
}

// Run returns the get result of the tracks.
func (uc *getAllTracksByAlbumIdUseCase) Run(ctx context.Context, id string) ([]*GetAllTracksByAlbumIdUseCaseOutputDto, error) {
	tracks, err := uc.trackRepo.FindByAlbumId(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}

	var getAllTracksByAlbumIdUseCaseOutputDtos []*GetAllTracksByAlbumIdUseCaseOutputDto
	for _, track := range tracks {
		artistNames := make([]string, len(track.Artists))
		for i, artist := range track.Artists {
			artistNames[i] = artist.Name
		}

		getAllTracksByArtistIdUseCaseOutputDto := &GetAllTracksByAlbumIdUseCaseOutputDto{
			ID:          track.ID.String(),
			Artists:     strings.Join(artistNames, ", "),
			Album:       track.Album.Name,
			Name:        track.Name,
			TrackNumber: track.TrackNumber,
			ReleaseDate: track.ReleaseDate,
		}
		getAllTracksByAlbumIdUseCaseOutputDtos = append(getAllTracksByAlbumIdUseCaseOutputDtos, getAllTracksByArtistIdUseCaseOutputDto)
	}

	return getAllTracksByAlbumIdUseCaseOutputDtos, nil
}
