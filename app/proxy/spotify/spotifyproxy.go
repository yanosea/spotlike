package spotifyproxy

import (
	"github.com/zmb3/spotify/v2"

	"github.com/yanosea/spotlike/app/proxy/http"
)

// Spotify is an interface for Spotify.
type Spotify interface {
	NewClient(httpClient *httpproxy.ClientInstance) ClientInstanceInterface
}

// SpotifyProxy is a struct that implements Spotify.
type SpotifyProxy struct{}

// New is a constructor for SpotifyProxy.
func New() Spotify {
	return &SpotifyProxy{}
}

func (*SpotifyProxy) NewClient(httpClient *httpproxy.ClientInstance) ClientInstanceInterface {
	return &ClientInstance{FieldClient: spotify.New(httpClient.FieldClient)}
}
