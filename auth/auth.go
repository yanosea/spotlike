package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/yanosea/spotlike/util"

	// https://github.com/fatih/color
	"github.com/fatih/color"
	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/thanhpk/randstr"
	"github.com/thanhpk/randstr"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
	// https://github.com/zmb3/spotify/v2/auth
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	auth_input_label_spotify_id                = "🆔 Input your Spotify Client ID"
	auth_input_label_spotify_secret            = "🔑 Input your Spotify Client Secret"
	auth_input_label_spotify_redirect_uri      = "🔗 Input your Spotify Redirect URI"
	auth_input_label_spotify_refresh_token     = "🔢 Input your Spotify Refresh Token if you have one (if you don't have it, leave it empty and press enter.)"
	auth_message_login_spotify                 = "🌐 Log in to Spotify by visiting the page below in your browser."
	auth_message_auth_success                  = "🎉 Authentication succeeded!"
	auth_message_suggest_set_env               = "If you don't want spotlike to ask questions above again, execute commands below to set envs or set your profile to set those."
	auth_message_template_set_env_command      = "export %s="
	auth_error_message_invalid_uri             = "Invalid URI..."
	auth_error_message_auth_failure            = "Authentication failed..."
	auth_error_message_template_refresh_failed = `Refresh failed...
Please clear your Spotify environment variables and try again.`
)

func IsEnvsSet() bool {
	return os.Getenv(util.AUTH_ENV_SPOTIFY_ID) != "" &&
		os.Getenv(util.AUTH_ENV_SPOTIFY_SECRET) != "" &&
		os.Getenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI) != "" &&
		os.Getenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN) != ""
}

func SetAuthInfo() {
	// SPOTIFY_ID
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

	// SPOTIFY_SECRET
	if os.Getenv(util.AUTH_ENV_SPOTIFY_SECRET) == "" {
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

	// SPOTIFY_REDIRECT_URI
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

	// SPOTIFY_REFRESH_TOKEN
	if os.Getenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN) == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_refresh_token,
		}

		input, _ := prompt.Run()
		os.Setenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN, input)
	}
}

func Authenticate(o io.Writer) (*spotify.Client, error) {
	port, err := getPortFromUri(os.Getenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI))
	if err != nil {
		return nil, err
	}

	var (
		client  *spotify.Client
		state   = randstr.Hex(11)
		channel = make(chan *spotify.Client)
	)

	authenticator := spotifyauth.New(
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserFollowModify,
			spotifyauth.ScopeUserLibraryModify,
		),
		spotifyauth.WithClientID(os.Getenv(util.AUTH_ENV_SPOTIFY_ID)),
		spotifyauth.WithClientSecret(os.Getenv(util.AUTH_ENV_SPOTIFY_SECRET)),
		spotifyauth.WithRedirectURL(os.Getenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI)),
	)

	if os.Getenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN) != "" {
		return refresh(authenticator, os.Getenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN))
	}

	url := authenticator.AuthURL(state)

	util.PrintWithWriterBetweenBlankLine(o, auth_message_login_spotify)
	util.PrintWithWriterWithBlankLineBelow(o, util.FormatIndent(url))

	var refreshToken string
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := authenticator.Token(r.Context(), state, r)

		if err != nil {
			http.Error(w, auth_error_message_auth_failure, http.StatusForbidden)
			return
		}

		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			return
		}

		client := spotify.New(authenticator.Client(r.Context(), tok))

		refreshToken = tok.RefreshToken
		channel <- client
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	go func() error {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			return err
		}
		return nil
	}()
	client = <-channel

	// print success message and suggest to set env
	util.PrintlnWithWriter(o, color.GreenString(auth_message_auth_success))
	util.PrintWithWriterBetweenBlankLine(o, color.YellowString(auth_message_suggest_set_env))
	util.PrintlnWithWriter(o, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_ID)+os.Getenv(util.AUTH_ENV_SPOTIFY_ID)))
	util.PrintlnWithWriter(o, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_SECRET)+os.Getenv(util.AUTH_ENV_SPOTIFY_SECRET)))
	util.PrintlnWithWriter(o, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_REDIRECT_URI)+os.Getenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI)))
	util.PrintlnWithWriter(o, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN)+refreshToken))

	return client, nil
}

func refresh(authenticator *spotifyauth.Authenticator, refreshToken string) (*spotify.Client, error) {
	tok := &oauth2.Token{
		TokenType:    "bearer",
		RefreshToken: refreshToken,
	}

	client := spotify.New(authenticator.Client(context.Background(), tok))
	if client == nil {
		clearCommands := []string{
			util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_ID)),
			util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_SECRET)),
			util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_REDIRECT_URI)),
			util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN)),
		}
		auth_error_message_refresh_failed := fmt.Sprintf("%s\n\n%s", auth_error_message_template_refresh_failed, strings.Join(clearCommands, "\n"))
		return nil, errors.New(auth_error_message_refresh_failed)
	}

	return client, nil
}

func getPortFromUri(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if u.Port() == "" {
		return "", errors.New(auth_error_message_invalid_uri)
	}

	return u.Port(), nil
}
