package spotlike

import (
	"os"

	c "github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/config"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/presenter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// AuthOptions provides the options for the auth command.
type AuthOptions struct {
	// ID is the client ID of the Spotify API.
	ID string
	// Secret is the client secret of the Spotify API.
	Secret string
	// RedirectUri is the redirect URI of the Spotify API.
	RedirectUri string
}

// NewAuthCommand returns a new instance of the auth command.
func NewAuthCommand(
	exit func(int),
	cobra proxy.Cobra,
	version string,
	conf *config.SpotlikeCliConfig,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("auth")
	cmd.SetAliases([]string{"au", "a"})
	cmd.SetUsageTemplate(authUsageTemplate)
	cmd.SetHelpTemplate(authHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(
		func(cmd *c.Command, _ []string) error {
			return runAuth(exit, cmd, version, conf, output)
		},
	)

	return cmd
}

// runAuth runs the auth command.
func runAuth(exit func(int), cmd *c.Command, version string, conf *config.SpotlikeCliConfig, output *string) error {
	clientManager := api.GetClientManager()
	if clientManager == nil {
		o := formatter.Red("❌ Client manager is not initialized...")
		*output = o
		return nil
	}
	if client, err := clientManager.GetClient(); err != nil && err.Error() != "client not initialized" {
		return err
	} else if client != nil {
		o := formatter.Yellow("⚡You've already setup your Spotify client...")
		*output = o
		return nil
	}

	if conf.SpotifyID == "" {
		for {
			if id, err := presenter.RunPrompt(
				"🆔 Input your Spotify Client ID",
			); err != nil && err.Error() == "^C" {
				if err := presenter.Print(os.Stdout, "\n"); err != nil {
					return err
				}
				if err := presenter.Print(os.Stdout, formatter.Yellow("🚫 Cancelled authentication...")); err != nil {
					return err
				}
				exit(130)
				return nil
			} else if err != nil {
				return err
			} else if id == "" {
				continue
			} else {
				conf.SpotifyID = id
				break
			}
		}
	}
	if conf.SpotifySecret == "" {
		for {
			if secret, err := presenter.RunPromptWithMask(
				"🔑 Input your Spotify Client Secret",
				'*',
			); err != nil && err.Error() == "^C" {
				if err := presenter.Print(os.Stdout, "\n"); err != nil {
					return err
				}
				if err := presenter.Print(os.Stdout, formatter.Yellow("🚫 Cancelled authentication...")); err != nil {
					return err
				}
				exit(130)
				return nil
			} else if err != nil {
				return err
			} else if secret == "" {
				continue
			} else {
				conf.SpotifySecret = secret
				break
			}
		}
	}
	if conf.SpotifyRedirectUri == "" {
		for {
			if redirectUri, err := presenter.RunPrompt(
				"🔗 Input your Spotify Redirect URI",
			); err != nil && err.Error() == "^C" {
				if err := presenter.Print(os.Stdout, "\n"); err != nil {
					return err
				}
				if err := presenter.Print(os.Stdout, formatter.Yellow("🚫 Cancelled authentication...")); err != nil {
					return err
				}
				exit(130)
				return nil
			} else if err != nil {
				return err
			} else if redirectUri == "" {
				continue
			} else {
				conf.SpotifyRedirectUri = redirectUri
				break
			}
		}
	}

	var (
		client       api.Client
		refreshToken string
		err          error
		authUrlChan  = make(chan string, 1)
		doneChan     = make(chan bool)
	)

	go func() {
		defer close(doneChan)
		err = clientManager.InitializeClient(
			cmd.Context(),
			&api.ClientConfig{
				SpotifyID:           conf.SpotifyID,
				SpotifySecret:       conf.SpotifySecret,
				SpotifyRedirectUri:  conf.SpotifyRedirectUri,
				SpotifyRefreshToken: "",
			},
		)
		if err != nil {
			authUrlChan <- ""
			return
		}

		client, err = clientManager.GetClient()
		if err != nil {
			authUrlChan <- ""
			return
		}

		_, refreshToken, err = client.Auth(authUrlChan, version)
		if err != nil {
			return
		}
		conf.SpotifyRefreshToken = refreshToken

		if client != nil {
			client.UpdateConfig(&api.ClientConfig{
				SpotifyID:           conf.SpotifyID,
				SpotifySecret:       conf.SpotifySecret,
				SpotifyRedirectUri:  conf.SpotifyRedirectUri,
				SpotifyRefreshToken: refreshToken,
			})
		}
	}()
	authUrl := <-authUrlChan

	if err := presenter.Print(os.Stdout, "\n"); err != nil {
		return err
	}
	if err := presenter.Print(os.Stdout, "🌐 Login to Spotify by visiting the page below in your browser."); err != nil {
		return err
	}
	if err := presenter.Print(os.Stdout, "  "+authUrl); err != nil {
		return err
	}

	<-doneChan
	if err != nil {
		o := formatter.Red("❌ Failed to authenticate... Please try again...")
		*output = o
		return err
	}

	if err := presenter.Print(os.Stdout, "\n"); err != nil {
		return err
	}
	if err := presenter.Print(os.Stdout, formatter.Green("🎉 Authentication succeeded!")); err != nil {
		return err
	}
	if err := presenter.Print(os.Stdout, formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.")); err != nil {
		return err
	}
	if err := presenter.Print(os.Stdout, "  export SPOTIFY_ID="+conf.SpotifyID); err != nil {
		return err
	}
	if err := presenter.Print(os.Stdout, "  export SPOTIFY_SECRET="+conf.SpotifySecret); err != nil {
		return err
	}
	if err := presenter.Print(os.Stdout, "  export SPOTIFY_REDIRECT_URI="+conf.SpotifyRedirectUri); err != nil {
		return err
	}
	if err := presenter.Print(os.Stdout, "  export SPOTIFY_REFRESH_TOKEN="+refreshToken); err != nil {
		return err
	}

	return nil
}

const (
	// authHelpTemplate is the help template of the auth command.
	authHelpTemplate = `🔑 Authenticate your Spotify client

You have to authenticate your Spotify client to use spotlike at first.
spotlike will ask you to input your Client ID, Client Secret, and Redirect URI.

Also, you can set those by environment variables below:
  SPOTIFY_ID
  SPOTIFY_SECRET
  SPOTIFY_REDIRECT_URI

Finally, spotlike will guide you to set those all environment variables.

` + authUsageTemplate
	// authUsageTemplate is the usage template of the auth command.
	authUsageTemplate = `Usage:
  spotlike auth [flags]
  spotlike au   [flags]
  spotlike a    [flags]

Flags:
  -h, --help           🤝 help for auth
`
)
