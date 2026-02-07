package repository

import (
	"context"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/zmb3/spotify/v2"
)

// artistRepository is a struct that implements the ArtistRepository interface.
type artistRepository struct {
	clientManager api.ClientManager
}

// NewArtistRepository returns a new instance of the artistRepository struct.
func NewArtistRepository() artistDomain.ArtistRepository {
	return &artistRepository{
		clientManager: api.GetClientManager(),
	}
}

// FindById returns the artist by the ID.
func (r *artistRepository) FindById(ctx context.Context, id spotify.ID) (*artistDomain.Artist, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	artist, err := client.GetArtist(ctx, id)
	if err != nil {
		return nil, err
	}

	return artistDomain.NewArtist(
		artist.ID,
		artist.Name,
	), nil
}

// FindByNameLimit returns the artist by the name with the limit.
func (r *artistRepository) FindByNameLimit(ctx context.Context, name string, limit int) ([]*artistDomain.Artist, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	result, err := client.Search(ctx, name, spotify.SearchTypeArtist, spotify.Limit(limit))
	if err != nil {
		return nil, err
	}

	var artists []*artistDomain.Artist
	for _, artist := range result.Artists.Artists {
		artists = append(
			artists,
			artistDomain.NewArtist(
				artist.ID,
				artist.Name,
			),
		)
	}

	return artists, nil
}

// IsLiked returns whether the artist is liked.
func (r *artistRepository) IsLiked(ctx context.Context, id spotify.ID) (bool, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return false, err
	}

	client := c.Open()
	result, err := client.CurrentUserFollows(ctx, "artist", id)
	if err != nil {
		return false, err
	}

	return result[0], nil
}

// Like likes the artist.
func (r *artistRepository) Like(ctx context.Context, id spotify.ID) error {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return err
	}

	client := c.Open()
	return client.FollowArtist(ctx, id)
}

// Unlike unlikes the artist.
func (r *artistRepository) Unlike(ctx context.Context, id spotify.ID) error {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return err
	}

	client := c.Open()
	return client.UnfollowArtist(ctx, id)
}
