package spotlike

import (
	"context"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
)

// getTrackUseCase is a struct that contains the use case of getting for an track.
type getTrackUseCase struct {
	trackRepo trackDomain.TrackRepository
}

// NewGetTrackUseCase returns a new instance of the GetTrackUseCase struct.
func NewGetTrackUseCase(trackRepo trackDomain.TrackRepository) *getTrackUseCase {
	return &getTrackUseCase{
		trackRepo: trackRepo,
	}
}

// GetTrackUseCaseOutputDto is a DTO struct that contains the output data of the getTrackUseCase.
type GetTrackUseCaseOutputDto struct {
	ID          string
	Name        string
	Artists     string
	Album       string
	TrackNumber spotify.Numeric
	ReleaseDate time.Time
}

// Run returns the get result of the track.
func (uc *getTrackUseCase) Run(ctx context.Context, id string) (*GetTrackUseCaseOutputDto, error) {
	track, err := uc.trackRepo.FindById(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}
	artistNames := make([]string, len(track.Artists))
	for i, artist := range track.Artists {
		artistNames[i] = artist.Name
	}

	return &GetTrackUseCaseOutputDto{
		ID:          track.ID.String(),
		Name:        track.Name,
		Artists:     strings.Join(artistNames, ", "),
		Album:       track.Album.Name,
		TrackNumber: track.TrackNumber,
		ReleaseDate: track.ReleaseDate,
	}, nil
}
