package album

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

// AlbumRepository is an interface that provides the repository for the album on Spotify.
type AlbumRepository interface {
	FindByArtistId(ctx context.Context, id spotify.ID) ([]*Album, error)
	FindById(ctx context.Context, id spotify.ID) (*Album, error)
	FindByNameLimit(ctx context.Context, name string, limit int) ([]*Album, error)
	IsLiked(ctx context.Context, id spotify.ID) (bool, error)
	Like(ctx context.Context, id spotify.ID) error
	Unlike(ctx context.Context, id spotify.ID) error
}
