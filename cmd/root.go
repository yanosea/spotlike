/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	"github.com/yanosea/spotlike/config"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

type Cmd struct {
	// version : the version of spotlike set from goreleaser
	version string
	// cfg : config of spotlike
	cfg *config.Config
	// rootCmd : root command of cobra
	rootCmd *cobra.Command
}

var (
	version = ""
)

func New() *Cmd {
	return &Cmd{
		version: version,
		cfg:     nil,
		rootCmd: &cobra.Command{
			Use:     "spotlike",
			Version: version,
			Short:   "'spotlike' is a CLI tool to LIKE contents in Spotify.",
			Long: `'spotlike' is a CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`,
		},
	}
}

func (cmd *Cmd) Init() error {
	// load config
	if err := cmd.loadConfig(); err != nil {
		return err
	}
	// chech tokens

	return nil
}

func (cmd *Cmd) loadConfig() error {
	var err error
	cmd.cfg, err = config.New()
	if err != nil {
		return err
	}
	return nil
}

func (cmd *Cmd) Execute() error {
	return cmd.rootCmd.Execute()
}
