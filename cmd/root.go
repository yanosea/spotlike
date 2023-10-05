/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/yanosea/spotlike/app"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

// version : the version of spotlike set from goreleaser
var version = ""

var rootCmd = &cobra.Command{
	Use:     "spotlike",
	Version: version,
	Short:   "'spotlike' is a CLI tool to LIKE contents in Spotify.",
	Long: `'spotlike' is a CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Init() {
	shared := app.New()

	// load cache
	cache, err := app.LoadCache()

	// check err
	if err != nil {
		os.Exit(1)
	}

	// chech tokens
	if cache.AccessToken == "" || cache.RefreshToken == "" {
	}
}
