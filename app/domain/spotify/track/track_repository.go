package track

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

// TrackRepository is an interface that provides the repository for the track on Spotify.
type TrackRepository interface {
	FindByAlbumId(ctx context.Context, id spotify.ID) ([]*Track, error)
	FindByArtistId(ctx context.Context, id spotify.ID) ([]*Track, error)
	FindById(ctx context.Context, id spotify.ID) (*Track, error)
	FindByNameLimit(ctx context.Context, name string, limit int) ([]*Track, error)
	IsLiked(ctx context.Context, id spotify.ID) (bool, error)
	Like(ctx context.Context, id spotify.ID) error
	Unlike(ctx context.Context, id spotify.ID) error
}
