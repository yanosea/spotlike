package main

import (
	"context"
	"os"

	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"
)

// SpotlikeCliParams is a struct that represents the options of spotlike cli.
type SpotlikeCliParams struct {
	// Version is the version of spotlike cli.
	Version string
	// Context is the context of spotlike cli.
	Context context.Context
	// Cobra is a proxy of spf13/cobra.
	Cobra proxy.Cobra
	// Envconfig is a proxy of kelseyhightower/envconfig.
	Envconfig proxy.Envconfig
	// Spotify is a proxy of zmb3/spotify.
	Spotify proxy.Spotify
	// Http is a proxy of net/http.
	Http proxy.Http
	// Randstr is a proxy of thanhpk/randstr
	Randstr proxy.Randstr
	// Url is a proxy of net/url.
	Url proxy.Url
	// VersionUtil provides the version of the application.
	VersionUtil utility.VersionUtil
}

var (
	// version is the version of spotlike cli and is embedded by goreleaser.
	version = ""
	// exit is a variable that contains the os.Exit function for injecting dependencies in testing.
	exit = os.Exit
	// spotlikeCliParams is a variable that contains the SpotlikeCliParams struct.
	spotlikeCliParams = SpotlikeCliParams{
		Version:     version,
		Context:     context.Background(),
		Cobra:       proxy.NewCobra(),
		Envconfig:   proxy.NewEnvconfig(),
		Spotify:     proxy.NewSpotify(),
		Http:        proxy.NewHttp(),
		Randstr:     proxy.NewRandstr(),
		Url:         proxy.NewUrl(),
		VersionUtil: utility.NewVersionUtil(proxy.NewDebug()),
	}
)

// main is the entry point of spotlike cli.
func main() {
	cli := command.NewCli(
		exit,
		spotlikeCliParams.Cobra,
		spotlikeCliParams.Context,
	)

	if exitCode := cli.Init(
		spotlikeCliParams.Envconfig,
		spotlikeCliParams.Spotify,
		spotlikeCliParams.Http,
		spotlikeCliParams.Randstr,
		spotlikeCliParams.Url,
		spotlikeCliParams.Version,
		spotlikeCliParams.VersionUtil,
	); exitCode != 0 {
		exit(exitCode)
	}

	defer spotlikeCliParams.Context.Done()
	exit(cli.Run())
}
