package cmd

import (
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

var version = "develop"

const (
	root_help_template = `⚪ spotlike is the CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.

Usage:
  spotlike [flags]
  spotlike [command]

Available Commands:
  auth        🔑 Authenticate your Spotify client.
  completion  🔧 Generate the autocompletion script for the specified shell.
  help        🤝 Help about any command
  like        🤍 Like content in Spotify by ID.
  search      🔍 Search for the ID of content in Spotify.
  version     🔖 Show the version of spotlike.

Flags:
  -h, --help      🤝 help for spotlike
  -v, --version   🔖 version for spotlike

Use "spotlike [command] --help" for more information about a command.
`
	root_use   = "spotlike"
	root_short = "⚪ spotlike is the CLI tool to LIKE contents in Spotify."
	root_long  = `⚪ spotlike is the CLI tool to LIKE contents in Spotify.

You can get the ID of some contents in Spotify.
You can LIKE the contents in Spotify by ID.`
	root_message_no_sub_command = `Use sub command below...

  * 🔑 auth
  * 🤍 like
    * 🎤  artist
    * 💿 album
    * 🎵 track
  * 🔍 search
  * 🔖 version

Use "spotlike --help" for more information about spotlike.`
)

type GlobalOption struct {
	Client *spotify.Client

	Out    io.Writer
	ErrOut io.Writer
}

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
			// If no sub command is specified, print the message and return nil.
			util.PrintlnWithWriter(o.Out, root_message_no_sub_command)

			return nil
		},
	}

	o.Out = outWriter
	o.ErrOut = errWriter
	cmd.SetOut(outWriter)
	cmd.SetErr(errWriter)

	cmd.SetHelpTemplate(root_help_template)

	cmd.AddCommand(
		newAuthCommand(o),
		newCompletionCommand(o),
		newLikeCommand(o),
		newSearchCommand(o),
		newVersionCommand(o),
	)

	return cmd, nil
}
