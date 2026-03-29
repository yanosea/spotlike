package like

import (
	c "github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// NewLikeCommand creates a new like command.
func NewLikeCommand(
	exit func(int),
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("like")
	cmd.SetAliases([]string{"li", "l"})
	cmd.SetUsageTemplate(likeUsageTemplate)
	cmd.SetHelpTemplate(likeHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.AddCommand(
		NewLikeArtistCommand(
			exit,
			cobra,
			authCmd,
			output,
		),
		NewLikeAlbumCommand(
			exit,
			cobra,
			authCmd,
			output,
		),
		NewLikeTrackCommand(
			exit,
			cobra,
			authCmd,
			output,
		),
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runLike(output)
		},
	)

	return cmd
}

// runLike runs the like command.
func runLike(output *string) error {
	o := formatter.Yellow("⚡ Use sub command below...")
	o += `

  - 🎤 artist
  - 💿 album
  - 🎵 track

Use "spotlike like --help" for more information about spotlike like.
Use "spotlike like [command] --help" for more information about a command.
`
	*output = o

	return nil
}

const (
	// likeHelpTemplate is the help template of the like command.
	likeHelpTemplate = `🤍 Like content on Spotify by ID.

You can like content on Spotify by ID.

Before using this command,
you need to get the ID of the content you want to like by using the search command.

` + likeUsageTemplate
	// likeUsageTemplate is the usage template of the like command.
	likeUsageTemplate = `Usage:
  spotlike like [flags]
  spotlike li   [flags]
  spotlike l    [flags]
  spotlike like [command]
  spotlike li   [command]
  spotlike l    [command]

Available Commands:
  artist, ar, A  🎤 Like an artist on Spotify by ID.
  album,  al, a  💿 Like albums on Spotify by ID.
  track,  tr, t  🎵 Like tracks on Spotify by ID.

Flags:
  -h, --help  🤝 help for like

Use "spotlike like [command] --help" for more information about a command.
`
)
