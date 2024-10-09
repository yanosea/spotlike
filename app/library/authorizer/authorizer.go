package authorizer

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/yanosea/spotlike/app/proxy/http"
	"github.com/yanosea/spotlike/app/proxy/os"
	"github.com/yanosea/spotlike/app/proxy/randstr"
	"github.com/yanosea/spotlike/app/proxy/spotify"
	"github.com/yanosea/spotlike/app/proxy/spotifyauth"
	"github.com/yanosea/spotlike/app/proxy/url"
)

// Authorizable is an interface that provides a function to authenticate with Spotify.
type Authorizable interface {
	Authenticate() (spotifyproxy.ClientInstanceInterface, AuthenticateStatus, error)
	GetAuthUrl() string
	IsEnvsSet() bool
	Refresh() (spotifyproxy.ClientInstanceInterface, AuthenticateStatus, error)
	SetAuthInfo(spotifyId string, spotifySecret string, spotifyRedirectUri string, spotifyRefreshToken string) SetEnvStatus
}

// Authorizer is a struct that implements Authorizable interface.
type Authorizer struct {
	Authenticator spotifyauthproxy.AuthenticatorInstanceInterface
	Http          httpproxy.Http
	Os            osproxy.Os
	Randstr       randstrproxy.Randstr
	RefreshToken  string
	Spotify       spotifyproxy.Spotify
	Spotifyauth   spotifyauthproxy.Spotifyauth
	State         string
	Url           urlproxy.Url
}

// New is a constructor of Authorizer.
func New(
	httpProxy httpproxy.Http,
	osProxy osproxy.Os,
	randstrProxy randstrproxy.Randstr,
	spotifyProxy spotifyproxy.Spotify,
	spotifyauthProxy spotifyauthproxy.Spotifyauth,
	urlProxy urlproxy.Url,
) *Authorizer {
	authorizer := &Authorizer{
		Authenticator: nil,
		Http:          httpProxy,
		Os:            osProxy,
		Randstr:       randstrProxy,
		RefreshToken:  "",
		Spotify:       spotifyProxy,
		Spotifyauth:   spotifyauthProxy,
		State:         "",
		Url:           urlProxy,
	}
	authorizer.Authenticator = authorizer.getAuthenticator()
	return authorizer
}

func (a *Authorizer) Authenticate() (*spotifyproxy.ClientInstanceInterface, AuthenticateStatus, error) {
	// get port from uri
	port, err := a.getPortFromUri(a.Os.Getenv(ENV_SPOTIFY_REDIRECT_URI))
	if err != nil {
		return nil, AuthenticateFailedInvalidUri, err
	}

	var (
		client  *spotifyproxy.ClientInstance
		channel = make(chan *spotifyproxy.ClientInstance)
	)

	a.Http.HandleFunc("/callback", func(w httpproxy.ResponseWriterInstanceInterface, r *httpproxy.RequestInstance) {
		tok, err := a.Authenticator.Token(r.Context(), a.State, r)
		if err != nil {
			a.Http.Error(w, err, http.StatusForbidden)
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

	return client, nil
}

func (a *Authorizer) GetAuthUrl() string {
	a.State = a.Randstr.Hex(11)
	return a.Authenticator.AuthURL(a.State)
}

// IsEnvsSet checks if the environment variables are set.
func (a *Authorizer) IsEnvsSet() bool {
	return a.Os.Getenv(ENV_SPOTIFY_ID) != "" &&
		a.Os.Getenv(ENV_SPOTIFY_SECRET) != "" &&
		a.Os.Getenv(ENV_SPOTIFY_REDIRECT_URI) != "" &&
		a.Os.Getenv(ENV_SPOTIFY_REFRESH_TOKEN) != ""
}

func (a *Authorizer) SetAuthInfo(
	spotifyId string,
	spotifySecret string,
	spotifyRedirectUri string,
	spotifyRefreshToken string,
) SetEnvStatus {
	// set SPOTIFY_ID
	if spotifyId != "" {
		a.Os.Setenv(ENV_SPOTIFY_ID, spotifyId)
	} else {
		return SetEnvFailedId
	}
	// set SPOTIFY_SECRET
	if spotifySecret != "" {
		a.Os.Setenv(ENV_SPOTIFY_SECRET, spotifySecret)
	} else {
		return SetEnvFailedSecret
	}
	// set SPOTIFY_REDIRECT_URI
	if spotifyRedirectUri != "" {
		a.Os.Setenv(ENV_SPOTIFY_REDIRECT_URI, spotifyRedirectUri)
	} else {
		return SetEnvFailedRedirectUri
	}
	// set SPOTIFY_REFRESH_TOKEN
	if spotifyRefreshToken != "" {
		a.Os.Setenv(ENV_SPOTIFY_REFRESH_TOKEN, spotifyRefreshToken)
	}

	return SetEnvSuccessfully
}

func (a *Authorizer) Refresh(refreshToken string) (*spotify.Client, error) {
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

// getAuthenticator gets authenticator.
func (a *Authorizer) getAuthenticator() spotifyauthproxy.AuthenticatorInstanceInterface {
	return a.Spotifyauth.NewAuthenticator(
		spotifyauthproxy.WithScopes(
			spotifyauthproxy.ScopeUserFollowRead,
			spotifyauthproxy.ScopeUserFollowModify,
			spotifyauthproxy.ScopeUserLibraryRead,
			spotifyauthproxy.ScopeUserLibraryModify,
		),
		spotifyauthproxy.WithClientID(a.Os.Getenv(ENV_SPOTIFY_ID)),
		spotifyauthproxy.WithClientSecret(a.Os.Getenv(ENV_SPOTIFY_SECRET)),
		spotifyauthproxy.WithRedirectURL(a.Os.Getenv(ENV_SPOTIFY_REDIRECT_URI)),
	)
}

// getPortFromUri gets port from uri.
func (a *Authorizer) getPortFromUri(uri string) (string, error) {
	u, err := a.Url.Parse(uri)
	if err != nil {
		return "", err
	}

	return u.Port(), nil
}
