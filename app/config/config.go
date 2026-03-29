package config

import (
	"github.com/yanosea/spotlike/pkg/proxy"
)

// Configurator is an interface that gets the configuration.
type Configurator interface {
	GetConfig() (*SpotlikeConfig, error)
}

// BaseConfigurator is a struct that implements the Configurator interface.
type BaseConfigurator struct {
	Envconfig proxy.Envconfig
}

// SpotlikeConfig is a struct that contains the configuration of the spotlike application.
type SpotlikeConfig struct {
	// SpotifyID is the Spotify client ID.
	SpotifyID string
	// SpotifySecret is the Spotify client secret.
	SpotifySecret string
	// SpotifyRedirectUri is the Spotify redirect URI.
	SpotifyRedirectUri string
	// SpotifyRefreshToken is the Spotify refresh token.
	SpotifyRefreshToken string
}

// NewConfigurator creates a new Configurator.
func NewConfigurator(
	envconfigProxy proxy.Envconfig,
) *BaseConfigurator {
	return &BaseConfigurator{
		Envconfig: envconfigProxy,
	}
}
