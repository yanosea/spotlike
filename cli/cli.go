package cli

import (
	"github.com/yanosea/spotlike/api"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// variables
var (
	// Client holds client of Spotify
	Client *spotify.Client
)

// Init sets the Spotify client to the global variable
func Init() error {
	// get the Spotify client
	client, err := api.GetClient()
	if err != nil {
		return err
	}

	// set the client
	Client = client
	return nil
}
