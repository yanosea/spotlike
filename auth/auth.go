package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"

	"github.com/yanosea/spotlike/util"

	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
	// https://github.com/zmb3/spotify/v2/auth
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	auth_env_spotify_id                                 = "SPOTIFY_ID"
	auth_env_spotify_secret                             = "SPOTIFY_SECRET"
	auth_env_spotify_redirect_uri                       = "SPOTIFY_REDIRECT_URI"
	auth_env_spotify_refresh_token                      = "SPOTIFY_REFRESH_TOKEN"
	auth_input_label_spotify_id                         = "Input your Spotify Client ID"
	auth_input_label_spotify_secret                     = "Input your Spotify Client Secret"
	auth_input_label_spotify_redirect_uri               = "Input your Spotify Redirect URI"
	auth_input_label_spotify_refresh_token              = "Input your Spotify Refresh Token if you have one (if you don't have it, leave it empty and press enter.)"
	auth_message_login_spotify                          = "Log in to Spotify by visiting the page below in your browser."
	auth_message_auth_success                           = "Authentication succeeded!"
	auth_message_suggest_set_env                        = "If you don't want spotlike to ask questions above again, execute commands below to set envs or set your profile to set those."
	auth_message_template_set_env_command               = "export %s="
	auth_error_message_auth_failure                     = "Authentication failed..."
	auth_error_message_invalid_uri                      = "Invalid URI..."
	auth_error_message_invalid_length_for_random_string = "Invalid length..."
)

var (
	authenticator *spotifyauth.Authenticator
	channel       = make(chan *spotify.Client)
	state, _      = generateRandomString(16)
)

func GetClient() (*spotify.Client, error) {
	setAuthInfo()

	client, err := authenticate()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func setAuthInfo() {
	if id := os.Getenv(auth_env_spotify_id); id == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_id,
		}

		input, _ := prompt.Run()
		os.Setenv(auth_env_spotify_id, input)
	}

	if secret := os.Getenv(auth_env_spotify_secret); secret == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_secret,
			Mask:  '*',
		}

		input, _ := prompt.Run()
		os.Setenv(auth_env_spotify_secret, input)
	}

	if uri := os.Getenv(auth_env_spotify_redirect_uri); uri == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_redirect_uri,
		}

		input, _ := prompt.Run()
		os.Setenv(auth_env_spotify_redirect_uri, input)
	}

	if refresh := os.Getenv(auth_env_spotify_refresh_token); refresh == "" {
		prompt := promptui.Prompt{
			Label: auth_input_label_spotify_refresh_token,
		}

		input, _ := prompt.Run()
		os.Setenv(auth_env_spotify_refresh_token, input)
	}

	authenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv(auth_env_spotify_redirect_uri)),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserFollowModify,
			spotifyauth.ScopeUserLibraryModify,
		),
	)
}

func authenticate() (*spotify.Client, error) {
	var client *spotify.Client

	refreshToken := os.Getenv(auth_env_spotify_refresh_token)
	if refreshToken == "" {
		port, err := getPortFromUri(os.Getenv(auth_env_spotify_redirect_uri))
		if err != nil {
			return nil, err
		}

		http.HandleFunc("/callback", completeAuthenticate)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
		go func() error {
			err := http.ListenAndServe(":"+port, nil)
			if err != nil {
				return err
			}
			return err
		}()
		url := authenticator.AuthURL(state)

		o := os.Stdout
		util.PrintWithWriterBetweenBlankLine(o, auth_message_login_spotify)
		util.PrintWithWriterWithBlankLineBelow(o, util.FormatIndent(url))

		client = <-channel
	} else {
		tok := &oauth2.Token{
			TokenType:    "bearer",
			RefreshToken: refreshToken,
		}

		client = spotify.New(authenticator.Client(context.Background(), tok))
	}

	return client, nil
}

func completeAuthenticate(w http.ResponseWriter, r *http.Request) {
	tok, err := authenticator.Token(r.Context(), state, r)

	if err != nil {
		http.Error(w, auth_error_message_auth_failure, http.StatusForbidden)
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
	}

	client := spotify.New(authenticator.Client(r.Context(), tok))

	o := os.Stdout
	util.PrintlnWithWriter(o, auth_message_auth_success)
	util.PrintWithWriterWithBlankLineBelow(o, auth_message_suggest_set_env)
	util.PrintlnWithWriter(o, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, auth_env_spotify_id)+os.Getenv(auth_env_spotify_id)))
	util.PrintlnWithWriter(o, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, auth_env_spotify_secret)+os.Getenv(auth_env_spotify_secret)))
	util.PrintlnWithWriter(o, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, auth_env_spotify_redirect_uri)+os.Getenv(auth_env_spotify_redirect_uri)))
	util.PrintWithWriterWithBlankLineBelow(o, util.FormatIndent(fmt.Sprintf(auth_message_template_set_env_command, auth_env_spotify_refresh_token)+tok.RefreshToken))

	channel <- client
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

// generateRandomString generates a random string of the specified length.
func generateRandomString(length int) (string, error) {
	if length < 0 {
		return "", errors.New(auth_error_message_invalid_length_for_random_string)
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}
