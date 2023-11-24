package cli

import (
	"github.com/yanosea/spotlike/api"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// variables
var (
	// version is version of spotlike
	Version = ""
	// Client is client of Spotify
	Client *spotify.Client
)

// SetClient sets the Spotify client to the global variable
func SetClient() error {
	// get the Spotify client
	if client, err := api.GetClient(); err != nil {
		return err
	} else {
		// set the client
		Client = client
	}
	return nil
}
