package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"
)

// checkLikeArtistUseCase is a struct that contains the use case of checking for an artist.
type checkLikeArtistUseCase struct {
	artistRepo artistDomain.ArtistRepository
}

// NewCheckLikeArtistUseCase returns a new instance of the checkLikeArtistUseCase struct.
func NewCheckLikeArtistUseCase(artistRepo artistDomain.ArtistRepository) *checkLikeArtistUseCase {
	return &checkLikeArtistUseCase{
		artistRepo: artistRepo,
	}
}

// Run returns the check result of the artist.
func (uc *checkLikeArtistUseCase) Run(ctx context.Context, id string) (bool, error) {
	return uc.artistRepo.IsLiked(ctx, spotify.ID(id))
}
