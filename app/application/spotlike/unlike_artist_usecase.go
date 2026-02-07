package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"
)

// unlikeArtistUseCase is a struct that contains the use case of unlikeing for an artist.
type unlikeArtistUseCase struct {
	artistRepo artistDomain.ArtistRepository
}

// NewUnlikeArtistUseCase returns a new instance of the LikeArtistUseCase struct.
func NewUnlikeArtistUseCase(artistRepo artistDomain.ArtistRepository) *unlikeArtistUseCase {
	return &unlikeArtistUseCase{
		artistRepo: artistRepo,
	}
}

// Run returns the unlike result of the artist.
func (uc *unlikeArtistUseCase) Run(ctx context.Context, id string) error {
	return uc.artistRepo.Unlike(ctx, spotify.ID(id))
}
