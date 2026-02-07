package command

import (
	"context"
	"os"

	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/config"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/presenter"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"
)

var (
	// output is the output string.
	output = ""
	// NewCli is a variable holding the current Cli creation function.
	NewCli CreateCliFunc = newCli
)

// Cli is an interface that represents the command line interface of spotlike cli.
type Cli interface {
	Init(envconfig proxy.Envconfig, spotify proxy.Spotify, http proxy.Http, randstr proxy.Randstr, url proxy.Url, version string, versionUtil utility.VersionUtil) int
	Run() int
}

// cli is a struct that represents the command line interface of spotlike cli.
type cli struct {
	Exit          func(int)
	Cobra         proxy.Cobra
	RootCommand   proxy.Command
	Context       context.Context
	ClientManager api.ClientManager
}

// CreateCliFunc is a function type for creating new Cli instances.
type CreateCliFunc func(exit func(int), cobra proxy.Cobra, ctx context.Context) Cli

// newCli is the default implementation of CreateCliFunc.
func newCli(exit func(int), cobra proxy.Cobra, ctx context.Context) Cli {
	return &cli{
		Exit:          exit,
		Cobra:         cobra,
		RootCommand:   nil,
		Context:       ctx,
		ClientManager: nil,
	}
}

// Init initializes the command line interface of spotlike.
func (c *cli) Init(
	envconfig proxy.Envconfig,
	spotify proxy.Spotify,
	http proxy.Http,
	randstr proxy.Randstr,
	url proxy.Url,
	version string,
	versionUtil utility.VersionUtil,
) int {
	configurator := config.NewSpotlikeCliConfigurator(envconfig)
	conf, err := configurator.GetConfig()
	if err != nil {
		output = formatter.AppendErrorToOutput(err, output)
		if err := presenter.Print(os.Stderr, output); err != nil {
			return 1
		}
		return 1
	}

	if c.ClientManager == nil {
		c.ClientManager = api.NewClientManager(spotify, http, randstr, url)
	}

	if conf.SpotifyID != "" &&
		conf.SpotifySecret != "" &&
		conf.SpotifyRedirectUri != "" &&
		conf.SpotifyRefreshToken != "" {
		if err := c.ClientManager.InitializeClient(
			c.Context,
			&api.ClientConfig{
				SpotifyID:           conf.SpotifyID,
				SpotifySecret:       conf.SpotifySecret,
				SpotifyRedirectUri:  conf.SpotifyRedirectUri,
				SpotifyRefreshToken: conf.SpotifyRefreshToken,
			},
		); err != nil {
			output = formatter.AppendErrorToOutput(err, output)
			if err := presenter.Print(os.Stderr, output); err != nil {
				return 1
			}
			return 1
		}
	}

	c.RootCommand = NewRootCommand(
		c.Exit,
		c.Cobra,
		versionUtil.GetVersion(version),
		conf,
		&output,
	)

	return 0
}

// Run runs the command line interface of spotlike cli.
func (c *cli) Run() (exitCode int) {
	exitCode = 0
	defer func() {
		if c.ClientManager != nil && c.ClientManager.IsClientInitialized() {
			if err := c.ClientManager.CloseClient(); err != nil {
				output = formatter.AppendErrorToOutput(err, output)
				if err := presenter.Print(os.Stderr, output); err != nil {
					exitCode = 1
				}
				exitCode = 1
			}
		}
	}()

	out := os.Stdout
	if err := c.RootCommand.ExecuteContext(c.Context); err != nil {
		output = formatter.AppendErrorToOutput(err, output)
		out = os.Stderr
		exitCode = 1
	}

	if output != "" {
		if err := presenter.Print(out, output); err != nil {
			exitCode = 1
		}
	}

	return exitCode
}
