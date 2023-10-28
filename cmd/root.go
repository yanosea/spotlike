/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

// version : version of spotlike
var version = ""

// rootCmd : cobra root command
var rootCmd = &cobra.Command{
	Use:     "spotlike",
	Version: version,
	Short:   "'spotlike' is a CLI tool to LIKE contents in Spotify.",
	Long: `'spotlike' is a CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`,
}

// Execute : entry point of root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

// init : executed before func Execute executed
func init() {
}
