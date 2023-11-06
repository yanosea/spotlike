package api

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"os"

	"github.com/yanosea/spotlike/app"

	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
	// https://github.com/zmb3/spotify/v2/auth
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// variables
var (
	// authenticator is Spotify authenticator
	authenticator *spotifyauth.Authenticator
	// channel is Spotify client
	channel = make(chan *spotify.Client)
	// state is Spotify auth state
	state = app.GenerateRandomString(16)
)

// GetClient returns a Spotify client
func GetClient() (*spotify.Client, error) {
	// get client info
	if err := setAuthInfo(); err != nil {
		return nil, err
	} else {
		// authenticate and get a spotify client
		if client, err := authenticate(); err != nil {
			return nil, err
		} else {
			return client, nil
		}
	}
}

// setAuthInfo sets Spotify authentication info to each environment variable
func setAuthInfo() error {
	// SPOTIFY_CLIENT_ID
	if id := os.Getenv("SPOTIFY_ID"); id == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Client ID",
		}

		input, err := prompt.Run()
		if err != nil {
			return err
		} else {
			os.Setenv("SPOTIFY_ID", input)
		}
	}

	// SPOTIFY_CLIENT_SECRET
	if secret := os.Getenv("SPOTIFY_SECRET"); secret == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Client Secret",
			Mask:  '*',
		}

		input, err := prompt.Run()
		if err != nil {
			return err
		} else {
			os.Setenv("SPOTIFY_SECRET", input)
		}
	}

	// SPOTIFY_REDIRECT_URI
	if uri := os.Getenv("SPOTIFY_REDIRECT_URI"); uri == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Redirect URI",
		}

		input, err := prompt.Run()
		if err != nil {
			return err
		} else {
			os.Setenv("SPOTIFY_REDIRECT_URI", input)
		}
	}

	// set authenticator
	authenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv("SPOTIFY_REDIRECT_URI")),
		spotifyauth.WithScopes(
			// to check the user has been already liked the artist
			spotifyauth.ScopeUserFollowRead,
			// to check the user has been already liked the album and the track
			spotifyauth.ScopeUserLibraryRead,
			// to like the artist
			spotifyauth.ScopeUserFollowModify,
			// to like the album and the track
			spotifyauth.ScopeUserLibraryModify),
	)

	// SPOTIFY_REFRESH_TOKEN
	if refresh := os.Getenv("SPOTIFY_REFRESH_TOKEN"); refresh == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Refresh Token if you have one (if you don't have one, leave it empty and press enter.)",
		}

		input, err := prompt.Run()
		if err != nil {
			return err
		} else {
			os.Setenv("SPOTIFY_REFRESH_TOKEN", input)
		}
	}
	return nil
}

// authenticate authenticates the auth info and returns a Spotify client
func authenticate() (*spotify.Client, error) {
	var client *spotify.Client

	// check the refresh token
	refreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")
	if refreshToken == "" {
		// get port from the redirect URI
		if portString, err := app.GetPortString(os.Getenv("SPOTIFY_REDIRECT_URI")); err != nil {
			return nil, err
		} else {
			// start an HTTP server
			http.HandleFunc("/callback", completeAuthenticate)
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
			go func() error {
				err := http.ListenAndServe(portString, nil)
				if err != nil {
					return err
				}
				return err
			}()

			// generate the Spotify authentication URI and print it
			url := authenticator.AuthURL(state)
			fmt.Printf("\nLog in to Spotify by visiting the page below in your browser.\n%s\n\n", url)

			// wait for authentication to complete and get a new Spotify client
			client = <-channel
		}
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

// completeAuthenticate completes the authentication process
func completeAuthenticate(w http.ResponseWriter, r *http.Request) {
	tok, err := authenticator.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
	}

	// get a new Spotify client
	client := spotify.New(authenticator.Client(r.Context(), tok))

	// print the refresh token and the suggestion message to set it to env
	fmt.Printf("Authentication succeeded!\n")
	fmt.Printf("If you don't want spotlike to ask questions above again, execute commands below to set envs or set your profile to set those.\n")
	fmt.Printf("  export SPOTIFY_ID=%s\n", os.Getenv("SPOTIFY_ID"))
	fmt.Printf("  export SPOTIFY_SECRET=%s\n", os.Getenv("SPOTIFY_SECRET"))
	fmt.Printf("  export SPOTIFY_REDIRECT_URI=%s\n", os.Getenv("SPOTIFY_REDIRECT_URI"))
	fmt.Printf("  export SPOTIFY_REFRESH_TOKEN=%s\n\n", tok.RefreshToken)
	channel <- client
}
