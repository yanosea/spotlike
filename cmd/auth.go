package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/yanosea/spotlike/auth"
	"github.com/yanosea/spotlike/util"

	// https://github.com/fatih/color
	"github.com/fatih/color"
	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	auth_use   = "auth"
	auth_short = "Authenticate your Spotify client."
	auth_long  = `Authenticate your Spotify client.

You have to authenticate your Spotify client to use spotlike at first.
spotlike will ask you to input your Client ID, Client Secret, Redirect URI, and Refresh Token.

You can use tha flags below.
The default is to save your credentials to the environment variables(e, env).
  * -c --cache (Save your credentials to the cache file.)
  * -e --env   (Save your credentials to the environment variables.)`
	auth_flag_cache                        = "cache"
	auth_shorthand_cache                   = "c"
	auth_flag_description_cache            = "save your credentials to the cache file"
	auth_flag_env                          = "env"
	auth_shorthand_env                     = "e"
	auth_flag_description_env              = "save your credentials to the environment variables."
	auth_input_label_spotify_id            = "Input your Spotify Client ID"
	auth_input_label_spotify_secret        = "Input your Spotify Client Secret"
	auth_input_label_spotify_redirect_uri  = "Input your Spotify Redirect URI"
	auth_input_label_spotify_refresh_token = "Input your Spotify Refresh Token if you have one (if you don't have it, leave it empty and press enter.)"
	auth_message_auth_success              = "Authentication succeeded!"
	auth_message_suggest_set_env           = "If you don't want spotlike to ask questions above again, execute commands below to set envs or set your profile to set those."
	auth_message_template_set_env_command  = "export %s="
)

type authOption struct {
	Cache bool
	Env   bool

	Out    io.Writer
	ErrOut io.Writer
}

func newAuthCommand(globalOption *GlobalOption) *cobra.Command {
	o := &authOption{}
	cmd := &cobra.Command{
		Use:   auth_use,
		Short: auth_short,
		Long:  auth_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return globalOption.auth()
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.Cache, auth_flag_cache, auth_shorthand_cache, false, auth_flag_description_cache)
	cmd.PersistentFlags().BoolVarP(&o.Env, auth_flag_env, auth_shorthand_env, true, auth_flag_description_env)
	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
}

func (g *GlobalOption) auth() error {
	setAuthInfo()

	_, refreshToken, err := auth.Authenticate()
	if err != nil {
		return err
	}

	util.PrintWithWriterWithBlankLineBelow(g.Out, color.GreenString(auth_message_auth_success))
	util.PrintWithWriterWithBlankLineBelow(g.Out, color.YellowString(auth_message_suggest_set_env))
	util.PrintlnWithWriter(g.Out, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_ID)+os.Getenv(util.AUTH_ENV_SPOTIFY_ID)))
	util.PrintlnWithWriter(g.Out, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_SECRET)+os.Getenv(util.AUTH_ENV_SPOTIFY_SECRET)))
	util.PrintlnWithWriter(g.Out, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_REDIRECT_URI)+os.Getenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI)))
	util.PrintWithWriterWithBlankLineBelow(g.Out, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN)+refreshToken))

	return nil
}

func setAuthInfo() {
	if os.Getenv(util.AUTH_ENV_SPOTIFY_ID) == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_id,
		}

		var input string
		for {
			input, _ = prompt.Run()
			if input != "" {
				break
			}
		}
		os.Setenv(util.AUTH_ENV_SPOTIFY_ID, input)
	}

	if spotifySecret := os.Getenv(util.AUTH_ENV_SPOTIFY_SECRET); spotifySecret == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_secret,
			Mask:  '*',
		}

		var input string
		for {
			input, _ = prompt.Run()
			if input != "" {
				break
			}
		}
		os.Setenv(util.AUTH_ENV_SPOTIFY_SECRET, input)
	}

	if os.Getenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI) == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_redirect_uri,
		}

		var input string
		for {
			input, _ = prompt.Run()
			if input != "" {
				break
			}
		}
		os.Setenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI, input)
	}

	if os.Getenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN) == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_refresh_token,
		}

		input, _ := prompt.Run()
		os.Setenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN, input)
	}
}
