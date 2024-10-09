package spotifyauthproxy

import (
	"github.com/zmb3/spotify/v2/auth"
)

// Spotifyauth is an interface for Spotifyauth.
type Spotifyauth interface {
	NewAuthenticator(options ...spotifyauth.AuthenticatorOption) *AuthenticatorInstance
}

// SpotifyauthProxy is a struct that implements Spotifyauth.
type SpotifyauthProxy struct{}

// New is a constructor for SpotifyauthProxy.
func New() Spotifyauth {
	return &SpotifyauthProxy{}
}

// NewAuthenticator is a proxy for spotifyauth.New().
func (p *SpotifyauthProxy) NewAuthenticator(options ...spotifyauth.AuthenticatorOption) *AuthenticatorInstance {
	return &AuthenticatorInstance{
		FieldAuthenticator: spotifyauth.New(options...),
	}
}
