package api

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"os"

	"github.com/yanosea/spotlike/util"

	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
	// https://github.com/zmb3/spotify/v2/auth
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// constants
const (
	// spotify_id is env string of the Spotify client ID.
	spotify_id = "SPOTIFY_ID"
	// input_label_spotify_id is the label of the Spotify client ID.
	input_label_spotify_id = "Input your Spotify Client ID"
	// spotify_secret is env string of the Spotify client secret.
	spotify_secret = "SPOTIFY_SECRET"
	// input_label_spotify_secret is the label of the Spotify client secret.
	input_label_spotify_secret = "Input your Spotify Client ID"
	// spotify_redirect_uri is env string of the Spotify redirect URI.
	spotify_redirect_uri = "SPOTIFY_REDIRECT_URI"
	// input_label_spotify_redirect_uri  is the label of the Spotify redirect URI.
	input_label_spotify_redirect_uri = "Input your Spotify Redirect URI"
	// spotify_refresh_token is env string of the Spotify refresh token.
	spotify_refresh_token = "SPOTIFY_REFRESH_TOKEN"
	// input_label_spotify_refresh_token is the label of the Spotify refresh token.
	input_label_spotify_refresh_token = "Input your Spotify Refresh Token if you have one (if you don't have it, leave it empty and press enter.)"
	// message_login_spotify is the message directing to login to Spotify
	message_login_spotify = "Log in to Spotify by visiting the page below in your browser."
	// error_message_auth_failure is the message for failure authenticating
	error_message_auth_failure = "Authentication failed..."
	// message_auth_success is the message for successful authenticating
	message_auth_success = "Authentication succeeded!"
	// message_suggest_set_env is the message for suggesting to set env variables
	message_suggest_set_env = "If you don't want spotlike to ask questions above again, execute commands below to set envs or set your profile to set those."
	// template_set_env_command is the template for suggesting to set env variables
	template_set_env_command = "export %s="
)

// variables
var (
	// authenticator is Spotify authenticator.
	authenticator *spotifyauth.Authenticator
	// channel is Spotify client.
	channel = make(chan *spotify.Client)
	// state is Spotify auth state.
	state, _ = generateRandomString(16)
)

// GetClient returns the Spotify client.
func GetClient() (*spotify.Client, error) {
	// get client info
	if err := setAuthInfo(); err != nil {
		return nil, err
	}

	// authenticate and get a spotify client
	client, err := authenticate()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// setAuthInfo sets Spotify authentication info to each environment variable.
func setAuthInfo() error {
	// SPOTIFY_CLIENT_ID
	if id := os.Getenv(spotify_id); id == "" {
		prompt := promptui.Prompt{
			Label: input_label_spotify_id,
		}

		input, err := prompt.Run()
		if err != nil {
			return err
		}

		os.Setenv(spotify_id, input)
	}

	// SPOTIFY_CLIENT_SECRET
	if secret := os.Getenv(spotify_secret); secret == "" {
		prompt := promptui.Prompt{
			Label: input_label_spotify_secret,
			Mask:  '*',
		}

		input, err := prompt.Run()
		if err != nil {
			return err
		}

		os.Setenv(spotify_secret, input)
	}

	// SPOTIFY_REDIRECT_URI
	if uri := os.Getenv(spotify_redirect_uri); uri == "" {
		prompt := promptui.Prompt{
			Label: input_label_spotify_redirect_uri,
		}

		input, err := prompt.Run()
		if err != nil {
			return err
		}

		os.Setenv(spotify_redirect_uri, input)
	}

	// set authenticator
	authenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv(spotify_redirect_uri)),
		spotifyauth.WithScopes(
			// to check the user has been already liked the artist
			spotifyauth.ScopeUserFollowRead,
			// to check the user has been already liked the album and the track
			spotifyauth.ScopeUserLibraryRead,
			// to like the artist
			spotifyauth.ScopeUserFollowModify,
			// to like the album and the track
			spotifyauth.ScopeUserLibraryModify,
		),
	)

	// SPOTIFY_REFRESH_TOKEN
	if refresh := os.Getenv(spotify_refresh_token); refresh == "" {
		prompt := promptui.Prompt{
			Label: input_label_spotify_refresh_token,
		}

		input, err := prompt.Run()
		if err != nil {
			return err
		}

		os.Setenv(spotify_refresh_token, input)
	}

	return nil
}

// authenticate authenticates the auth info and returns a Spotify client.
func authenticate() (*spotify.Client, error) {
	var client *spotify.Client

	// check the refresh token
	refreshToken := os.Getenv(spotify_refresh_token)
	if refreshToken == "" {
		// get port from the redirect URI
		port, err := getPortFromUri(os.Getenv(spotify_redirect_uri))
		if err != nil {
			return nil, err
		}

		// start an HTTP server
		http.HandleFunc("/callback", completeAuthenticate)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
		go func() error {
			err := http.ListenAndServe(port, nil)
			if err != nil {
				return err
			}
			return err
		}()

		// generate the Spotify authentication URI and print it
		url := authenticator.AuthURL(state)
		util.PrintBetweenBlankLine(message_login_spotify)
		util.PrintlnWithBlankLineBelow(util.FormatIndent(url))

		// wait for authentication to complete and get a new Spotify client
		client = <-channel
	} else {
		// refresh token
		tok := &oauth2.Token{
			TokenType:    "bearer",
			RefreshToken: refreshToken,
		}

		// get a new Spotify client
		client = spotify.New(authenticator.Client(context.Background(), tok))
	}

	return client, nil
}

// completeAuthenticate completes the authentication process.
func completeAuthenticate(w http.ResponseWriter, r *http.Request) {
	tok, err := authenticator.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, error_message_auth_failure, http.StatusForbidden)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
	}

	// get a new Spotify client
	client := spotify.New(authenticator.Client(r.Context(), tok))

	// print the the suggestion message to set env variables
	fmt.Println(message_auth_success)
	util.PrintlnWithBlankLineBelow(message_suggest_set_env)
	fmt.Println(util.FormatIndent(fmt.Sprintf(template_set_env_command, spotify_id) + os.Getenv(spotify_id)))
	fmt.Println(util.FormatIndent(fmt.Sprintf(template_set_env_command, spotify_secret) + os.Getenv(spotify_secret)))
	fmt.Println(util.FormatIndent(fmt.Sprintf(template_set_env_command, spotify_redirect_uri) + os.Getenv(spotify_redirect_uri)))
	util.PrintlnWithBlankLineBelow(util.FormatIndent(fmt.Sprintf(template_set_env_command, spotify_refresh_token) + tok.RefreshToken))

	channel <- client
}
