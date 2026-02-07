package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"
)

// likeArtistUseCase is a struct that contains the use case of likeing for an artist.
type likeArtistUseCase struct {
	artistRepo artistDomain.ArtistRepository
}

// NewLikeArtistUseCase returns a new instance of the LikeArtistUseCase struct.
func NewLikeArtistUseCase(artistRepo artistDomain.ArtistRepository) *likeArtistUseCase {
	return &likeArtistUseCase{
		artistRepo: artistRepo,
	}
}

// Run returns the like result of the artist.
func (uc *likeArtistUseCase) Run(ctx context.Context, id string) error {
	return uc.artistRepo.Like(ctx, spotify.ID(id))
}
