package authorizer

import (
	"github.com/yanosea/spotlike/app/proxy/context"
	"github.com/yanosea/spotlike/app/proxy/http"
	"github.com/yanosea/spotlike/app/proxy/oauth2"
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
	Refresh() (spotifyproxy.ClientInstanceInterface, AuthenticateStatus)
	SetAuthInfo(spotifyId string, spotifySecret string, spotifyRedirectUri string, spotifyRefreshToken string) SetEnvStatus
}

// Authorizer is a struct that implements Authorizable interface.
type Authorizer struct {
	Authenticator spotifyauthproxy.AuthenticatorInstanceInterface
	Context       contextproxy.Context
	Http          httpproxy.Http
	Oauth2        oauth2proxy.Oauth2
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
	contextProxy contextproxy.Context,
	httpProxy httpproxy.Http,
	oauth2Proxy oauth2proxy.Oauth2,
	osProxy osproxy.Os,
	randstrProxy randstrproxy.Randstr,
	spotifyProxy spotifyproxy.Spotify,
	spotifyauthProxy spotifyauthproxy.Spotifyauth,
	urlProxy urlproxy.Url,
) *Authorizer {
	authorizer := &Authorizer{
		Authenticator: nil,
		Context:       contextProxy,
		Http:          httpProxy,
		Oauth2:        oauth2Proxy,
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

func (a *Authorizer) Authenticate() (spotifyproxy.ClientInstanceInterface, AuthenticateStatus, error) {
	// get port from uri
	port, err := a.getPortFromUri(a.Os.Getenv(ENV_SPOTIFY_REDIRECT_URI))
	if err != nil {
		return nil, AuthenticateFailedInvalidUri, err
	}

	var (
		client  spotifyproxy.ClientInstanceInterface
		channel = make(chan spotifyproxy.ClientInstanceInterface)
	)

	a.Http.HandleFunc("/callback", func(w httpproxy.ResponseWriterInstanceInterface, r *httpproxy.RequestInstance) {
		tok, err := a.Authenticator.Token(r.Context(), a.State, r)
		if err != nil {
			a.Http.Error(w, err, httpproxy.StatusForbidden)
			return
		}

		if st := r.FormValue("state"); st != a.State {
			a.Http.NotFound(w, r)
			return
		}

		client := a.Spotify.NewClient(a.Authenticator.Client(r.Context(), tok))

		a.RefreshToken = tok.FieldToken.RefreshToken
		channel <- client
	})

	a.Http.HandleFunc("/", func(w httpproxy.ResponseWriterInstanceInterface, r *httpproxy.RequestInstance) {})
	go func() error {
		err := a.Http.ListenAndServe(":"+port, nil)
		if err != nil {
			return err
		}
		return nil
	}()
	client = <-channel

	if client == nil {
		return nil, AuthenticateFailed, nil
	}

	return client, AuthenticatedSuccessfully, nil
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

func (a *Authorizer) Refresh() (spotifyproxy.ClientInstanceInterface, AuthenticateStatus) {
	tok := a.Oauth2.NewToken("bearer", a.RefreshToken)
	if client := a.Spotify.NewClient(a.Authenticator.Client(a.Context.Background(), tok)); client != nil {
		return client, AuthenticatedSuccessfully
	} else {
		return nil, AuthenticateFailed
	}
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
