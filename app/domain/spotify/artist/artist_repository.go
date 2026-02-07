package artist

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

// ArtistRepository is an interface that provides the repository for the artist on Spotify.
type ArtistRepository interface {
	FindById(ctx context.Context, id spotify.ID) (*Artist, error)
	FindByNameLimit(ctx context.Context, name string, limit int) ([]*Artist, error)
	IsLiked(ctx context.Context, id spotify.ID) (bool, error)
	Like(ctx context.Context, id spotify.ID) error
	Unlike(ctx context.Context, id spotify.ID) error
}
