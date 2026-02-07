package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
)

// unlikeTrackUseCase is a struct that contains the use case of unlikeing for an track.
type unlikeTrackUseCase struct {
	trackDomain trackDomain.TrackRepository
}

// NewUnlikeTrackUseCase returns a new instance of the UnlikeTrackUseCase struct.
func NewUnlikeTrackUseCase(trackDomain trackDomain.TrackRepository) *unlikeTrackUseCase {
	return &unlikeTrackUseCase{
		trackDomain: trackDomain,
	}
}

// Run returns the unlike result of the artist.
func (uc *unlikeTrackUseCase) Run(ctx context.Context, id string) error {
	return uc.trackDomain.Unlike(ctx, spotify.ID(id))
}
