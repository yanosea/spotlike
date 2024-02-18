package cmd

import (
	"errors"
	"io"
	"os"

	"github.com/yanosea/spotlike/util"

	// https://github.com/fatih/color
	"github.com/fatih/color"
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

const (
	root_use   = "spotlike"
	root_short = "spotlike is the CLI tool to LIKE contents in Spotify."
	root_long  = `'spotlike' is the CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`
	root_error_message_no_sub_command = `Use sub command below...
  * auth
  * like
  * search
  * version`
)

type GlobalOption struct {
	Client *spotify.Client

	Out    io.Writer
	ErrOut io.Writer
}

var version = "develop"

func Execute() int {
	o := os.Stdout
	e := os.Stderr

	rootCmd, err := NewRootCommand(o, e)
	if err != nil {
		util.PrintlnWithWriter(e, color.RedString(err.Error()))
		return 1
	}

	if err = rootCmd.Execute(); err != nil {
		util.PrintlnWithWriter(e, color.RedString(err.Error()))
		return 1
	}

	return 0
}

func NewRootCommand(outWriter, errWriter io.Writer) (*cobra.Command, error) {
	o := &GlobalOption{}

	cmd := &cobra.Command{
		Use:           root_use,
		Short:         root_short,
		Long:          root_long,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New(root_error_message_no_sub_command)
		},
	}

	o.Out = outWriter
	o.ErrOut = errWriter
	cmd.SetOut(outWriter)
	cmd.SetErr(errWriter)

	cmd.AddCommand(
		newAuthCommand(o),
		newLikeCommand(o),
		newSearchCommand(o),
		newVersionCommand(o),
	)

	return cmd, nil
}
