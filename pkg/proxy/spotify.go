package proxy

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
)

// Spotify is an interface that provides a proxy of the methods of spotify.
type Spotify interface {
	NewAuthenticator(opts ...spotifyauth.AuthenticatorOption) Authenticator
	NewClient(client *http.Client, opts ...spotify.ClientOption) Client
}

// spotifyProxy is a proxy struct that implements the Spotify interface.
type spotifyProxy struct{}

// NewSpotify is a proxy method that returns the Spotify.
func NewSpotify() Spotify {
	return &spotifyProxy{}
}

// NewAuthenticator is a proxy method that returns the spotify.Authenticator.
func (*spotifyProxy) NewAuthenticator(opts ...spotifyauth.AuthenticatorOption) Authenticator {
	return &authenticatorProxy{authenticator: spotifyauth.New(opts...)}
}

// NewClient is a proxy method that returns the spotify.Client.
func (*spotifyProxy) NewClient(client *http.Client, opts ...spotify.ClientOption) Client {
	return &clientProxy{client: spotify.New(client, opts...)}
}

// Authenticator is an interface that provides a proxy of the methods of spotify.Authenticator.
type Authenticator interface {
	AuthURL(state string) string
	Client(ctx context.Context, tok *oauth2.Token) *http.Client
	Token(ctx context.Context, state string, r *http.Request) (*oauth2.Token, error)
}

// authenticatorProxy is a proxy struct that implements the Authenticator interface.
type authenticatorProxy struct {
	authenticator *spotifyauth.Authenticator
}

// AuthURL is a proxy method that calls the AuthURL method of the spotify.Authenticator.
func (a *authenticatorProxy) AuthURL(state string) string {
	return a.authenticator.AuthURL(state)
}

// Client is a proxy method that calls the Client method of the spotify.Authenticator.
func (a *authenticatorProxy) Client(ctx context.Context, tok *oauth2.Token) *http.Client {
	return a.authenticator.Client(ctx, tok)
}

// Token is a proxy method that calls the Token method of the spotify.Authenticator.
func (a *authenticatorProxy) Token(ctx context.Context, state string, r *http.Request) (*oauth2.Token, error) {
	return a.authenticator.Token(ctx, state, r)
}

// Client is an interface that provides a proxy of the methods of spotify.Client.
type Client interface {
	AddAlbumsToLibrary(ctx context.Context, ids ...spotify.ID) error
	AddTracksToLibrary(ctx context.Context, ids ...spotify.ID) error
	CurrentUserFollows(ctx context.Context, t string, ids ...spotify.ID) ([]bool, error)
	FollowArtist(ctx context.Context, id spotify.ID) error
	GetAlbum(ctx context.Context, id spotify.ID, opts ...spotify.RequestOption) (*spotify.FullAlbum, error)
	GetAlbumTracks(ctx context.Context, id spotify.ID, opts ...spotify.RequestOption) (*spotify.SimpleTrackPage, error)
	GetArtist(ctx context.Context, id spotify.ID) (*spotify.FullArtist, error)
	GetArtistAlbums(ctx context.Context, artistId spotify.ID, ts []spotify.AlbumType, opts ...spotify.RequestOption) (*spotify.SimpleAlbumPage, error)
	GetTrack(ctx context.Context, id spotify.ID, opts ...spotify.RequestOption) (*spotify.FullTrack, error)
	RemoveAlbumsFromLibrary(ctx context.Context, ids ...spotify.ID) error
	RemoveTracksFromLibrary(ctx context.Context, ids ...spotify.ID) error
	Search(ctx context.Context, query string, t spotify.SearchType, opts ...spotify.RequestOption) (*spotify.SearchResult, error)
	UnfollowArtist(ctx context.Context, id spotify.ID) error
	UserHasAlbums(ctx context.Context, ids ...spotify.ID) ([]bool, error)
	UserHasTracks(ctx context.Context, ids ...spotify.ID) ([]bool, error)
}

// clientProxy is a proxy struct that implements the Client interface.
type clientProxy struct {
	client *spotify.Client
}

// AddAlbumsToLibrary is a proxy method that calls the AddAlbumsToLibrary method of the spotify.Client.
func (c *clientProxy) AddAlbumsToLibrary(ctx context.Context, ids ...spotify.ID) error {
	return c.client.AddAlbumsToLibrary(ctx, ids...)
}

// AddTracksToLibrary is a proxy method that calls the AddTracksToLibrary method of the spotify.Client.
func (c *clientProxy) AddTracksToLibrary(ctx context.Context, ids ...spotify.ID) error {
	return c.client.AddTracksToLibrary(ctx, ids...)
}

// CurrentUserFollows is a proxy method that calls the CurrentUserFollows method of the spotify.Client.
func (c *clientProxy) CurrentUserFollows(ctx context.Context, t string, ids ...spotify.ID) ([]bool, error) {
	return c.client.CurrentUserFollows(ctx, t, ids...)
}

// FollowArtist is a proxy method that calls the FollowArtist method of the spotify.Client.
func (c *clientProxy) FollowArtist(ctx context.Context, id spotify.ID) error {
	return c.client.FollowArtist(ctx, id)
}

// GetAlbum is a proxy method that calls the GetAlbum method of the spotify.Client.
func (c *clientProxy) GetAlbum(ctx context.Context, id spotify.ID, opts ...spotify.RequestOption) (*spotify.FullAlbum, error) {
	return c.client.GetAlbum(ctx, id, opts...)
}

// GetAlbumTracks is a proxy method that calls the GetAlbumTracks method of the spotify.Client.
func (c *clientProxy) GetAlbumTracks(ctx context.Context, id spotify.ID, opts ...spotify.RequestOption) (*spotify.SimpleTrackPage, error) {
	return c.client.GetAlbumTracks(ctx, id, opts...)
}

// GetArtist is a proxy method that calls the GetArtist method of the spotify.Client.
func (c *clientProxy) GetArtist(ctx context.Context, id spotify.ID) (*spotify.FullArtist, error) {
	return c.client.GetArtist(ctx, id)
}

// GetArtistAlbums is a proxy method that calls the GetArtistAlbums method of the spotify.Client.
func (c *clientProxy) GetArtistAlbums(ctx context.Context, artistID spotify.ID, ts []spotify.AlbumType, opts ...spotify.RequestOption) (*spotify.SimpleAlbumPage, error) {
	return c.client.GetArtistAlbums(ctx, artistID, ts, opts...)
}

// GetTrack is a proxy method that calls the GetTrack method of the spotify.Client.
func (c *clientProxy) GetTrack(ctx context.Context, id spotify.ID, opts ...spotify.RequestOption) (*spotify.FullTrack, error) {
	return c.client.GetTrack(ctx, id, opts...)
}

// RemoveAlbumsFromLibrary is a proxy method that calls the RemoveAlbumsFromLibrary method of the spotify.Client.
func (c *clientProxy) RemoveAlbumsFromLibrary(ctx context.Context, ids ...spotify.ID) error {
	return c.client.RemoveAlbumsFromLibrary(ctx, ids...)
}

// RemoveTracksFromLibrary is a proxy method that calls the RemoveTracksFromLibrary method of the spotify.Client.
func (c *clientProxy) RemoveTracksFromLibrary(ctx context.Context, ids ...spotify.ID) error {
	return c.client.RemoveTracksFromLibrary(ctx, ids...)
}

// Search is a proxy method that calls the Search method of the spotify.Client.
func (c *clientProxy) Search(ctx context.Context, query string, t spotify.SearchType, opts ...spotify.RequestOption) (*spotify.SearchResult, error) {
	return c.client.Search(ctx, query, t, opts...)
}

// UnfollowArtist is a proxy method that calls the UnfollowArtist method of the spotify.Client.
func (c *clientProxy) UnfollowArtist(ctx context.Context, id spotify.ID) error {
	return c.client.UnfollowArtist(ctx, id)
}

// UserHasAlbums is a proxy method that calls the UserHasAlbums method of the spotify.Client.
func (c *clientProxy) UserHasAlbums(ctx context.Context, ids ...spotify.ID) ([]bool, error) {
	return c.client.UserHasAlbums(ctx, ids...)
}

// UserHasTracks is a proxy method that calls the UserHasTracks method of the spotify.Client.
func (c *clientProxy) UserHasTracks(ctx context.Context, ids ...spotify.ID) ([]bool, error) {
	return c.client.UserHasTracks(ctx, ids...)
}
