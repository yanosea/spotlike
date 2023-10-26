/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

var version = ""

var rootCmd = &cobra.Command{
	Use:     "spotlike",
	Version: version,
	Short:   "'spotlike' is a CLI tool to LIKE contents in Spotify.",
	Long: `'spotlike' is a CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`,
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
}
