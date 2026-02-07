package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
)

// checkLikeAlbumUseCase is a struct that contains the use case of checking for an album.
type checkLikeAlbumUseCase struct {
	albumRepo albumDomain.AlbumRepository
}

// NewCheckLikeAlbumUseCase returns a new instance of the checkLikeAlbumUseCase struct.
func NewCheckLikeAlbumUseCase(albumRepo albumDomain.AlbumRepository) *checkLikeAlbumUseCase {
	return &checkLikeAlbumUseCase{
		albumRepo: albumRepo,
	}
}

// Run returns the check result of the album.
func (uc *checkLikeAlbumUseCase) Run(ctx context.Context, id string) (bool, error) {
	return uc.albumRepo.IsLiked(ctx, spotify.ID(id))
}
