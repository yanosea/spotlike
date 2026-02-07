package unlike

import (
	c "github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// NewUnlikeCommand creates a new unlike command.
func NewUnlikeCommand(
	exit func(int),
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("unlike")
	cmd.SetAliases([]string{"un", "u"})
	cmd.SetUsageTemplate(unlikeUsageTemplate)
	cmd.SetHelpTemplate(unlikeHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.AddCommand(
		NewUnlikeArtistCommand(
			exit,
			cobra,
			authCmd,
			output,
		),
		NewUnlikeAlbumCommand(
			exit,
			cobra,
			authCmd,
			output,
		),
		NewUnlikeTrackCommand(
			exit,
			cobra,
			authCmd,
			output,
		),
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runUnlike(output)
		},
	)

	return cmd
}

// runUnlike runs the unlike command.
func runUnlike(output *string) error {
	o := formatter.Yellow("⚡ Use sub command below...")
	o += `

  - 🎤 artist
  - 💿 album
  - 🎵 track

Use "spotlike unlike --help" for more information about spotlike unlike.
Use "spotlike unlike [command] --help" for more information about a command.
`
	*output = o

	return nil
}

const (
	// unlikeHelpTemplate is the help template of the unlike command.
	unlikeHelpTemplate = `💔 Unlike content on Spotify by ID.

You can unlike content on Spotify by ID.

Before using this command,
you need to get the ID of the content you want to like by using the search command.

` + unlikeUsageTemplate
	// unlikeUsageTemplate is the usage template of the unlike command.
	unlikeUsageTemplate = `Usage:
  spotlike unlike [flags]
  spotlike un     [flags]
  spotlike u      [flags]
  spotlike unlike [command]
  spotlike un     [command]
  spotlike u      [command]

Available Commands:
  artist, ar, A  🎤 Unlike an artist on Spotify by ID.
  album,  al, a  💿 Unlike albums on Spotify by ID.
  track,  tr, t  🎵 Unlike tracks on Spotify by ID.

Flags:
  -h, --help  🤝 help for unlike

Use "spotlike unlike [command] --help" for more information about a command.
`
)
