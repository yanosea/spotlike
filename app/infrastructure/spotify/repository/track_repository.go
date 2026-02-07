package repository

import (
	"context"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/zmb3/spotify/v2"
)

// trackRepository is a struct that implements the TrackRepository interface.
type trackRepository struct {
	clientManager api.ClientManager
}

// NewTrackRepository returns a new instance of the trackRepository struct.
func NewTrackRepository() trackDomain.TrackRepository {
	return &trackRepository{
		clientManager: api.GetClientManager(),
	}
}

// FindByArtistId returns the tracks by the artist ID.
func (r *trackRepository) FindByArtistId(ctx context.Context, id spotify.ID) ([]*trackDomain.Track, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	albumsResult, err := client.GetArtistAlbums(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	var tracks []*trackDomain.Track
	for _, album := range albumsResult.Albums {
		tracksResult, err := client.GetAlbumTracks(ctx, album.ID)
		if err != nil {
			return nil, err
		}

		for _, track := range tracksResult.Tracks {
			tracks = append(
				tracks,
				trackDomain.NewTrack(
					track.ID,
					track.Name,
					track.Artists,
					album,
					track.TrackNumber,
					album.ReleaseDateTime(),
				),
			)
		}
	}

	return tracks, nil
}

// FindByAlbumId returns the tracks by the album ID.
func (r *trackRepository) FindByAlbumId(ctx context.Context, id spotify.ID) ([]*trackDomain.Track, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	album, err := client.GetAlbum(ctx, id)
	if err != nil {
		return nil, err
	}

	var tracks []*trackDomain.Track
	tracksResult, err := client.GetAlbumTracks(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, track := range tracksResult.Tracks {
		tracks = append(
			tracks,
			trackDomain.NewTrack(
				track.ID,
				track.Name,
				track.Artists,
				album.SimpleAlbum,
				track.TrackNumber,
				album.ReleaseDateTime(),
			),
		)
	}

	return tracks, nil
}

// FindById returns the track by the ID.
func (r *trackRepository) FindById(ctx context.Context, id spotify.ID) (*trackDomain.Track, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	track, err := client.GetTrack(ctx, id)
	if err != nil {
		return nil, err
	}

	return trackDomain.NewTrack(
		track.ID,
		track.Name,
		track.Artists,
		track.Album,
		track.TrackNumber,
		track.Album.ReleaseDateTime(),
	), nil
}

// FindByNameLimit returns the track by the name with the limit.
func (r *trackRepository) FindByNameLimit(ctx context.Context, name string, limit int) ([]*trackDomain.Track, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return nil, err
	}

	client := c.Open()
	result, err := client.Search(ctx, name, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		return nil, err
	}

	var tracks []*trackDomain.Track
	for _, track := range result.Tracks.Tracks {
		tracks = append(
			tracks,
			trackDomain.NewTrack(
				track.ID,
				track.Name,
				track.Artists,
				track.Album,
				track.TrackNumber,
				track.Album.ReleaseDateTime(),
			),
		)
	}

	return tracks, nil
}

// IsLiked returns whether the track is liked.
func (r *trackRepository) IsLiked(ctx context.Context, id spotify.ID) (bool, error) {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return false, err
	}

	client := c.Open()
	result, err := client.UserHasTracks(ctx, id)
	if err != nil {
		return false, err
	}

	return result[0], nil
}

// Like likes the track.
func (r *trackRepository) Like(ctx context.Context, id spotify.ID) error {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return err
	}

	client := c.Open()
	return client.AddTracksToLibrary(ctx, id)
}

// Unlike unlikes the track.
func (r *trackRepository) Unlike(ctx context.Context, id spotify.ID) error {
	c, err := r.clientManager.GetClient()
	if err != nil {
		return err
	}

	client := c.Open()
	return client.RemoveTracksFromLibrary(ctx, id)
}
