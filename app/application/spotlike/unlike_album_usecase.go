package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
)

// unlikeAlbumUseCase is a struct that contains the use case of unlikeing for an album.
type unlikeAlbumUseCase struct {
	albumDomain albumDomain.AlbumRepository
}

// NewUnlikeAlbumUseCase returns a new instance of the unLikeAlbumUseCase struct.
func NewUnlikeAlbumUseCase(albumDomain albumDomain.AlbumRepository) *unlikeAlbumUseCase {
	return &unlikeAlbumUseCase{
		albumDomain: albumDomain,
	}
}

// Run returns the unlike result of the album.
func (uc *unlikeAlbumUseCase) Run(ctx context.Context, id string) error {
	return uc.albumDomain.Unlike(ctx, spotify.ID(id))
}
