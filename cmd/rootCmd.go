package cmd

import (
	"github.com/yanosea/spotlike/cli"
	"github.com/yanosea/spotlike/constants"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

// variables
var (
	// version is version of spotlike
	version = ""
)

// rootCmd is the Cobra root command of spotlike
var rootCmd = &cobra.Command{
	Use:     constants.RootUse,
	Version: version,
	Short:   constants.RootShort,
	Long:    constants.RootLong,
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
	// init the Spotify client
	if err := cli.Init(); err != nil {
		return err
	}

	return nil
}
