package command

import (
	c "github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command/spotlike"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command/spotlike/completion"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command/spotlike/get"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command/spotlike/like"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command/spotlike/unlike"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/config"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// RootOptions provides the options for the root command.
type RootOptions struct {
	// Version is a flag to show the version of spotlike.
	Version bool
}

var (
	// rootOps is a variable to store the root options with the default values for injecting the dependencies in testing.
	rootOps = RootOptions{
		Version: false,
	}
)

// NewRootCommand returns a new instance of the root command.
func NewRootCommand(
	exit func(int),
	cobra proxy.Cobra,
	version string,
	conf *config.SpotlikeCliConfig,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("spotlike")
	cmd.SetUsageTemplate(rootUsageTemplate)
	cmd.SetHelpTemplate(rootHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.PersistentFlags().BoolVarP(
		&rootOps.Version,
		"version",
		"v",
		false,
		"🔖 show the version of spotlike",
	)
	versionCmd := spotlike.NewVersionCommand(
		cobra,
		version,
		output,
	)
	authCmd := spotlike.NewAuthCommand(
		exit,
		cobra,
		version,
		conf,
		output,
	)
	cmd.AddCommand(
		authCmd,
		completion.NewCompletionCommand(
			cobra,
			output,
		),
		get.NewGetCommand(
			cobra,
			authCmd,
			output,
		),
		like.NewLikeCommand(
			exit,
			cobra,
			authCmd,
			output,
		),
		unlike.NewUnlikeCommand(
			exit,
			cobra,
			authCmd,
			output,
		),
		spotlike.NewSearchCommand(
			cobra,
			authCmd,
			output,
		),
		versionCmd,
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runRoot(cmd, args, output, versionCmd)
		},
	)

	return cmd
}

// runRoot runs the root command.
func runRoot(
	cmd *c.Command,
	args []string,
	output *string,
	versionCmd proxy.Command,
) error {
	if rootOps.Version {
		return versionCmd.RunE(cmd, args)
	}

	o := formatter.Yellow("⚡ Use sub command below...")
	o += `

- 🔑 auth,       au,   a - Authenticate your Spotify client.
- 📚 get,        ge,   g - Get the information of the content on Spotify by ID.
- 🤍 like,       li,   l - Like content on Spotify by ID.
- 💔 unlike,     un,   u - Unlike content on Spotify by ID.
- 🔍 search,     se,   s - Search for the ID of content in Spotify.
- 🔧 completion, comp, c - Generate the autocompletion script for the specified shell.
- 🔖 version,    ver,  v - Show the version of spotlike.
- 🤝 help                - Help for spotlike.

Use "spotlike --help" for more information about spotlike.
Use "spotlike [command] --help" for more information about a command.
`
	*output = o

	return nil
}

const (
	// rootHelpTemplate is the help template of the root command.
	rootHelpTemplate = `⚪ spotlike is the CLI tool to LIKE contents in Spotify.

You can like tracks, albums, and artists in Spotify by ID.
Also, you can search for the ID of tracks, albums, and artists in Spotify.

` + rootUsageTemplate
	// rootUsageTemplate is the usage template of the root command.
	rootUsageTemplate = `Usage:
  spotlike [flags]
  spotlike [command]

Available Commands:
  auth,       au,   a  🔑 Authenticate your Spotify client.
  get,        ge,   g  📚 Get the information of the content on Spotify by ID.
  like,       li,   l  🤍 Like content on Spotify by ID.
  unlike,     un,   u  💔 Unlike content on Spotify by ID.
  search,     se,   s  🔍 Search for the ID of content in Spotify.
  completion, comp, c  🔧 Generate the autocompletion script for the specified shell.
  version,    ver,  v  🔖 Show the version of spotlike.
  help                 🤝 Help for spotlike.

Flags:
  -h, --help     🤝 help for spotlike
  -v, --version  🔖 version for spotlike

Use "spotlike [command] --help" for more information about a command.
`
)
