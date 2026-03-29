package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
)

// likeTrackUseCase is a struct that contains the use case of likeing for an track.
type likeTrackUseCase struct {
	trackDomain trackDomain.TrackRepository
}

// NewLikeTrackUseCase returns a new instance of the LikeTrackUseCase struct.
func NewLikeTrackUseCase(trackDomain trackDomain.TrackRepository) *likeTrackUseCase {
	return &likeTrackUseCase{
		trackDomain: trackDomain,
	}
}

// Run returns the like result of the track.
func (uc *likeTrackUseCase) Run(ctx context.Context, id string) error {
	return uc.trackDomain.Like(ctx, spotify.ID(id))
}
