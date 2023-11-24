package cmd

import (
	"github.com/yanosea/spotlike/cli"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

// rootCmd is the Cobra root command of spotlike
var rootCmd = &cobra.Command{
	Use:     cli.ROOT.RootName(),
	Version: cli.Version,
	Short:   cli.ROOT.RootShort(),
	Long:    cli.ROOT.RootLong(),
}

// Execute is entry point of the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

// Init is executed before func Execute is executed
func Init() error {
	// set the Spotify client
	if err := cli.SetClient(); err != nil {
		return err
	}

	return nil
}
