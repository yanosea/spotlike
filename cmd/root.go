/*
Package cmd defines and configures the sub commands for the spotlike CLI application using Cobra.
*/
package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

// version of spotlike
var version = ""

// rootCmd is the Cobra root command of spotlike
var rootCmd = &cobra.Command{
	Use:     "spotlike",
	Version: version,
	Short:   "'spotlike' is a CLI tool to LIKE contents in Spotify.",
	Long: `'spotlike' is a CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`,
}

// Execute is entry point of the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

// init is executed before func Execute is executed
func init() {
}
