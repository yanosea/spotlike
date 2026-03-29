package api

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"

	"github.com/zmb3/spotify/v2/auth"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// Client is an interface that contains the client.
type Client interface {
	Auth(chan<- string, string) (proxy.Client, string, error)
	Close() error
	Open() proxy.Client
	UpdateConfig(*ClientConfig)
}

// client is a struct that contains the client.
type client struct {
	spotify proxy.Spotify
	client  proxy.Client
	http    proxy.Http
	randstr proxy.Randstr
	url     proxy.Url
	config  *ClientConfig
	context context.Context
	mutex   *sync.RWMutex
}

var (
	// dc is a channel that is used to notify the server shutdown.
	dc chan struct{}
)

// Close closes the client.
func (c *client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.client != nil {
		c.client = nil
	}

	return nil
}

// Open opens the client.
func (c *client) Open() proxy.Client {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.client != nil {
		return c.client
	}

	authenticator := c.spotify.NewAuthenticator(
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserFollowModify,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserLibraryModify,
		),
		spotifyauth.WithClientID(c.config.SpotifyID),
		spotifyauth.WithClientSecret(c.config.SpotifySecret),
		spotifyauth.WithRedirectURL(c.config.SpotifyRedirectUri),
	)

	tok := &oauth2.Token{
		TokenType:    "bearer",
		RefreshToken: c.config.SpotifyRefreshToken,
	}

	c.client = c.spotify.NewClient(authenticator.Client(c.context, tok))

	return c.client
}

// UpdateConfig updates the client configuration.
func (c *client) UpdateConfig(config *ClientConfig) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.config = config
}

// Auth authenticates the client.
func (c *client) Auth(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
	c.mutex.Lock()
	if c.client != nil {
		c.mutex.Unlock()
		return c.client, "", nil
	}
	c.mutex.Unlock()

	ctx, cancel := context.WithTimeout(c.context, 5*time.Minute)
	defer cancel()

	authenticator := c.spotify.NewAuthenticator(
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserFollowModify,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserLibraryModify,
		),
		spotifyauth.WithClientID(c.config.SpotifyID),
		spotifyauth.WithClientSecret(c.config.SpotifySecret),
		spotifyauth.WithRedirectURL(c.config.SpotifyRedirectUri),
	)

	state := c.randstr.Hex(11)
	authUrlChan <- authenticator.AuthURL(state)

	uri, err := c.url.Parse(c.config.SpotifyRedirectUri)
	if err != nil {
		return nil, "", err
	}
	var port string
	if port = uri.Port(); port == "" {
		return nil, "", errors.New("failed to get port")
	}
	server := c.http.NewServer(":" + port)

	var (
		client       proxy.Client
		refreshToken string
		clientChan   = make(chan proxy.Client, 1)
		errChan      = make(chan error, 1)
		doneChan     = make(chan struct{}, 1)
	)
	dc = doneChan

	defer func() {
		close(clientChan)
		close(errChan)
		close(doneChan)
	}()

	c.http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			select {
			case doneChan <- struct{}{}:
			default:
			}
		}()

		tok, err := authenticator.Token(r.Context(), state, r)
		if err != nil {
			errChan <- err
			http.Error(w, "authentication failed: "+err.Error(), http.StatusForbidden)
			return
		}

		if tok == nil || tok.RefreshToken == "" {
			errChan <- errors.New("refresh token is empty")
			http.Error(w, "authentication failed: refresh token is empty", http.StatusForbidden)
			return
		}

		client = c.spotify.NewClient(authenticator.Client(r.Context(), tok))
		refreshToken = tok.RefreshToken
		clientChan <- client

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		successPageHtml := `
<!DOCTYPE html>
<html>
  <head>
    <title>🎉 spotlike authentication succeeded!</title>
  </head>
  <body>
    <p>🎉 Authentication succeeded! You may close this page and return to spotlike!</p>
    <p style="font-style: italic">- <a href="https://github.com/yanosea/spotlike">spotlike</a> v` + version + `</p>
  </body>
</html>`
		if _, err := w.Write([]byte(successPageHtml)); err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	cleanup := func() error {
		ctx, cancel := context.WithTimeout(c.context, 1*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}

	select {
	case <-ctx.Done():
		if err := cleanup(); err != nil {
			return nil, "", err
		}
		return nil, "", errors.New("authentication timed out")
	case err := <-errChan:
		if cerr := cleanup(); cerr != nil {
			return nil, "", err
		}
		return nil, "", err
	case <-doneChan:
		if err := cleanup(); err != nil {
			return nil, "", err
		}

		select {
		case client = <-clientChan:
			if client == nil {
				return nil, "", errors.New("failed to get client")
			}
		default:
			return nil, "", errors.New("client not received")
		}
	}

	c.mutex.Lock()
	c.client = client
	c.mutex.Unlock()

	return client, refreshToken, nil
}
