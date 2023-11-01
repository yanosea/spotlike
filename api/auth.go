/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"regexp"

	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
	// https://github.com/zmb3/spotify/v2/auth
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	// auth : spotify Authenticator
	auth *spotifyauth.Authenticator
	// ch : spotify client
	ch = make(chan *spotify.Client)
	// state : spotify auth state
	state = generateRandomString(16)
)

// GetClient : returns spotify client
func GetClient() (*spotify.Client, error) {
	// get client info
	if err := setAuthInfo(); err != nil {
		return nil, err
	} else {
		// authenticate and get spotify client
		if client, err := authenticate(); err != nil {
			return nil, err
		} else {
			return client, nil
		}
	}
}

// setAuthInfo : set spotify authenticate info to each env
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

	// set auhenticator
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv("SPOTIFY_REDIRECT_URI")),
		spotifyauth.WithScopes(spotifyauth.ScopeUserFollowModify, spotifyauth.ScopeUserLibraryModify),
	)

	// SPOTIFY_REFRESH_TOKEN
	if refresh := os.Getenv("SPOTIFY_REFRESH_TOKEN"); refresh == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Refresh Token if you heve (if you don't have, leave it empty and press enter.)",
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

// authenticate : authenticates info and returns spotify client with context
func authenticate() (*spotify.Client, error) {
	var client *spotify.Client

	// check refresh token
	refreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")
	if refreshToken == "" {
		// get portString from redirect uri
		if portString, err := getPortString(os.Getenv("SPOTIFY_REDIRECT_URI")); err != nil {
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

			// generate spotify authenticate url and show it
			url := auth.AuthURL(state)
			fmt.Printf("Log in to Spotify by visiting the page below in your browser.\n%s\n\n", url)

			// wait for auth to complete
			client = <-ch
		}
	} else {
		// refresh token
		tok := &oauth2.Token{
			TokenType:    "bearer",
			RefreshToken: refreshToken,
		}

		// get new client
		client = spotify.New(auth.Client(context.Background(), tok))
	}

	return client, nil
}

// completeAuthenticate : complete authenticate
func completeAuthenticate(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Printf("Your Refresh Token :\t%s\n", tok.RefreshToken)
	fmt.Printf("Set this token to the env 'SPOTIFY_REFRESH_TOKEN'.\n\n")
	ch <- client
}

// getPortString : returns port string like ':xxxx'
func getPortString(uri string) (string, error) {
	pattern := `:(\d+)`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(uri)

	if len(match) > 1 {
		return match[0], nil
	} else {
		return "", errors.New("invalid redirect uri")
	}
}

// generateRandomString : generate random string
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
