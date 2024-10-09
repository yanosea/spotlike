package spotifyproxy

// Spotify is an interface for Spotify.
type Spotify interface {
}

// SpotifyProxy is a struct that implements Spotify.
type SpotifyProxy struct{}

// New is a constructor for SpotifyProxy.
func New() Spotify {
	return &SpotifyProxy{}
}
