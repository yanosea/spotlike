package config

import (
	baseConfig "github.com/yanosea/spotlike/app/config"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// SpotlikeCliConfigurator is an interface that gets the configuration of the Spotlike cli application.
type SpotlikeCliConfigurator interface {
	GetConfig() (*SpotlikeCliConfig, error)
}

// cliConfigurator is a struct that implements the SpotlikeCliConfigurator interface.
type cliConfigurator struct {
	*baseConfig.BaseConfigurator
}

// NewSpotlikeCliConfigurator creates a new SpotlikeCliConfigurator.
func NewSpotlikeCliConfigurator(
	envconfigProxy proxy.Envconfig,
) SpotlikeCliConfigurator {
	return &cliConfigurator{
		BaseConfigurator: baseConfig.NewConfigurator(
			envconfigProxy,
		),
	}
}

// SpotlikeCliConfig is a struct that contains the configuration of the Spotlike cli application.
type SpotlikeCliConfig struct {
	baseConfig.SpotlikeConfig
}

// envConfig is a struct that contains the environment variables.
type envConfig struct {
	// SpotifyID is the Spotify client ID.
	SpotifyID string `envconfig:"SPOTIFY_ID"`
	// SpotifySecret is the Spotify client secret.
	SpotifySecret string `envconfig:"SPOTIFY_SECRET"`
	// SpotifyRedirectUri is the Spotify redirect URI.
	SpotifyRedirectUri string `envconfig:"SPOTIFY_REDIRECT_URI"`
	// SpotifyRefreshToken is the Spotify refresh token.
	SpotifyRefreshToken string `envconfig:"SPOTIFY_REFRESH_TOKEN"`
}

// GetConfig gets the configuration of the spotlike cli application.
func (c *cliConfigurator) GetConfig() (*SpotlikeCliConfig, error) {
	var env envConfig
	if err := c.Envconfig.Process("", &env); err != nil {
		return nil, err
	}

	config := &SpotlikeCliConfig{
		SpotlikeConfig: baseConfig.SpotlikeConfig{
			SpotifyID:           env.SpotifyID,
			SpotifySecret:       env.SpotifySecret,
			SpotifyRedirectUri:  env.SpotifyRedirectUri,
			SpotifyRefreshToken: env.SpotifyRefreshToken,
		},
	}

	return config, nil
}
