package api

import ()

// ClientConfig is a struct that contains the configuration of the client.
type ClientConfig struct {
	// SpotifyID is the Spotify client ID.
	SpotifyID string
	// SpotifySecret is the Spotify client secret.
	SpotifySecret string
	// SpotifyRedirectUri is the Spotify redirect URI.
	SpotifyRedirectUri string
	// SpotifyRefreshToken is the Spotify refresh token.
	SpotifyRefreshToken string
}
