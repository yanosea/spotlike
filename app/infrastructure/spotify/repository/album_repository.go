package repository

import (
	"context"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/zmb3/spotify/v2"
)

// albumRepository is a struct that implements the AlbumRepository interface.
type albumRepository struct {
	clientManager api.ClientManager
}

// NewAlbumRepository returns a new instance of the albumRepository struct.
func NewAlbumRepository() albumDomain.AlbumRepository {
	return &albumRepository{
		clientManager: api.GetClientManager(),
	}
}

// FindByArtistId returns the albums by the artist ID.
func (r *albumRepository) FindByArtistId(ctx context.Context, id spotify.ID) ([]*albumDomain.Album, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	result, err := client.GetArtistAlbums(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	var albums []*albumDomain.Album
	for _, album := range result.Albums {
		albums = append(
			albums,
			albumDomain.NewAlbum(
				album.ID,
				album.Name,
				album.Artists,
				album.ReleaseDateTime(),
			),
		)
	}

	return albums, nil
}

// FindById returns the album by the ID.
func (r *albumRepository) FindById(ctx context.Context, id spotify.ID) (*albumDomain.Album, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	album, err := client.GetAlbum(ctx, id)
	if err != nil {
		return nil, err
	}

	return albumDomain.NewAlbum(
		album.ID,
		album.Name,
		album.Artists,
		album.ReleaseDateTime(),
	), nil
}

// FindByNameLimit returns the album by the name with the limit.
func (r *albumRepository) FindByNameLimit(ctx context.Context, name string, limit int) ([]*albumDomain.Album, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	result, err := client.Search(ctx, name, spotify.SearchTypeAlbum, spotify.Limit(limit))
	if err != nil {
		return nil, err
	}

	var albums []*albumDomain.Album
	for _, album := range result.Albums.Albums {
		albums = append(
			albums,
			albumDomain.NewAlbum(
				album.ID,
				album.Name,
				album.Artists,
				album.ReleaseDateTime(),
			),
		)
	}

	return albums, nil
}

// IsLiked returns whether the album is liked.
func (r *albumRepository) IsLiked(ctx context.Context, id spotify.ID) (bool, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return false, err
	}

	client := c.Open()
	result, err := client.UserHasAlbums(ctx, id)
	if err != nil {
		return false, err
	}

	return result[0], nil
}

// Like likes the album.
func (r *albumRepository) Like(ctx context.Context, id spotify.ID) error {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return err
	}

	client := c.Open()
	return client.AddAlbumsToLibrary(ctx, id)
}

// Unlike unlikes the album.
func (r *albumRepository) Unlike(ctx context.Context, id spotify.ID) error {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return err
	}

	client := c.Open()
	return client.RemoveAlbumsFromLibrary(ctx, id)
}
