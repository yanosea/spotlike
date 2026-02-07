package get

import (
	c "github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// NewGetCommand creates a new get command.
func NewGetCommand(
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("get")
	cmd.SetAliases([]string{"ge", "g"})
	cmd.SetUsageTemplate(getUsageTemplate)
	cmd.SetHelpTemplate(getHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.AddCommand(
		NewGetAlbumsCommand(
			cobra,
			authCmd,
			output,
		),
		NewGetTracksCommand(
			cobra,
			authCmd,
			output,
		),
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runGet(output)
		},
	)

	return cmd
}

// runGet runs the get command.
func runGet(output *string) error {
	o := formatter.Yellow("⚡ Use sub command below...")
	o += `

  - 💿 album
  - 🎵 track

Use "spotlike get --help" for more information about spotlike get.
Use "spotlike get [command] --help" for more information about a command.
`
	*output = o

	return nil
}

const (
	// getHelpTemplate is the help template of the get command.
	getHelpTemplate = `ℹ️  Get the information of the content on Spotify by ID.

You can get the information of the content on Spotify by ID.

Before using this command,
you need to get the ID of the content you want to get by using the search command.

` + getUsageTemplate
	// getUsageTemplate is the usage template of the get command.
	getUsageTemplate = `Usage:
  spotlike get [flags]
  spotlike ge  [flags]
  spotlike g   [flags]
  spotlike get [command]
  spotlike ge  [command]
  spotlike g   [command]

Available Commands:
  albums,  als, a  💿 Get the information of the albums on Spotify by ID.
  tracks,  trs, t  🎵 Get the information of the tracks on Spotify by ID.

Flags:
  -h, --help  🤝 help for get

Use "spotlike get [command] --help" for more information about a command.
`
)
