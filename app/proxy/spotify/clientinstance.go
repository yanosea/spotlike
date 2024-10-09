package spotifyproxy

import (
	"github.com/zmb3/spotify/v2"
)

// ClientInstanceInterface is an interface for spotify.Client.
type ClientInstanceInterface interface {
}

// ClientInstance is a struct that implements Client.
type ClientInstance struct {
	FieldClient *spotify.Client
}
