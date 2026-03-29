package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
)

// checkLikeTrackUseCase is a struct that contains the use case of checking for an track.
type checkLikeTrackUseCase struct {
	trackRepo trackDomain.TrackRepository
}

// NewCheckLikeTrackUseCase returns a new instance of the checkLikeTrackUseCase struct.
func NewCheckLikeTrackUseCase(trackRepo trackDomain.TrackRepository) *checkLikeTrackUseCase {
	return &checkLikeTrackUseCase{
		trackRepo: trackRepo,
	}
}

// Run returns the check result of the track.
func (uc *checkLikeTrackUseCase) Run(ctx context.Context, id string) (bool, error) {
	return uc.trackRepo.IsLiked(ctx, spotify.ID(id))
}
