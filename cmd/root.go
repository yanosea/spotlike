package cmd

import (
	"github.com/yanosea/spotlike/api"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// constants
const (
	// root_use is the useage of root command.
	root_use = "spotlike"
	// root_short is the root_short description of root command.
	root_short = "spotlike is the CLI tool to LIKE contents in Spotify."
	// root_long is the root_long description of root command.
	root_long = `'spotlike' is the CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`
)

// variables
var (
	// version is version of spotlike
	version = ""
	// Client holds client of Spotify
	Client *spotify.Client
)

// rootCmd is the Cobra root command of spotlike.
var rootCmd = &cobra.Command{
	Use:     root_use,
	Version: version,
	Short:   root_short,
	Long:    root_long,
}

// Init is executed before func Execute is executed.
func Init() error {
	client, err := api.GetClient()
	if err != nil {
		return err
	}

	// set the client
	Client = client

	return nil
}

// Execute is the entry point of the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
