package spotlike

import (
	"context"

	"github.com/zmb3/spotify/v2"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
)

// likeAlbumUseCase is a struct that contains the use case of likeing for an album.
type likeAlbumUseCase struct {
	albumDomain albumDomain.AlbumRepository
}

// NewLikeAlbumUseCase returns a new instance of the LikeAlbumUseCase struct.
func NewLikeAlbumUseCase(albumDomain albumDomain.AlbumRepository) *likeAlbumUseCase {
	return &likeAlbumUseCase{
		albumDomain: albumDomain,
	}
}

// Run returns the like result of the album.
func (uc *likeAlbumUseCase) Run(ctx context.Context, id string) error {
	return uc.albumDomain.Like(ctx, spotify.ID(id))
}
