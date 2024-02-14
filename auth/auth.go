package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"

	"github.com/yanosea/spotlike/util"

	// https://github.com/thanhpk/randstr"
	"github.com/thanhpk/randstr"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
	// https://github.com/zmb3/spotify/v2/auth
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	auth_close_tab_html                       = "<!DOCTYPE html><html><head><script>open(location, '_self').close();</script></head></html>"
	auth_error_message_auth_failure           = "Authentication failed..."
	auth_error_message_invalid_uri            = "Invalid URI..."
	auth_error_message_platform_not_supported = "Platform not supported..."
	auth_error_message_refresh_failed         = `Refresh failed...
  Please clear your Spotify environment variables and try again.`
	auth_message_login_spotify = "Log in to Spotify by visiting the page below in your browser."
)

func Authenticate() (*spotify.Client, string, error) {
	port, err := getPortFromUri(os.Getenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI))
	if err != nil {
		return nil, "", err
	}

	var (
		client  *spotify.Client
		state   = randstr.Hex(11)
		channel = make(chan *spotify.Client)
	)

	authenticator := spotifyauth.New(
		spotifyauth.WithRedirectURL(os.Getenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI)),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserFollowModify,
			spotifyauth.ScopeUserLibraryModify,
		),
		spotifyauth.WithClientID(os.Getenv(util.AUTH_ENV_SPOTIFY_ID)),
		spotifyauth.WithClientSecret(os.Getenv(util.AUTH_ENV_SPOTIFY_SECRET)),
	)

	if os.Getenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN) != "" {
		return refresh(authenticator, os.Getenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN))
	}

	url := authenticator.AuthURL(state)

	o := os.Stdout
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

	return client, refreshToken, nil
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

func refresh(authenticator *spotifyauth.Authenticator, refreshToken string) (*spotify.Client, string, error) {
	tok := &oauth2.Token{
		TokenType:    "bearer",
		RefreshToken: refreshToken,
	}

	client := spotify.New(authenticator.Client(context.Background(), tok))
	if client == nil {
		return nil, "", errors.New(auth_error_message_platform_not_supported)
	}

	return client, tok.RefreshToken, nil
}
