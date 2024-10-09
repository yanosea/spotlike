// authenticator.go
package spotifyauthproxy

import (
	"github.com/zmb3/spotify/v2/auth"

	"github.com/yanosea/spotlike/app/proxy/context"
	"github.com/yanosea/spotlike/app/proxy/http"
	"github.com/yanosea/spotlike/app/proxy/oauth2"
)

// AuthenticatorInstanceInterface is an interface for AuthenticatorInstance.
type AuthenticatorInstanceInterface interface {
	AuthURL(state string) string
	Token(ctx *contextproxy.ContextInstance, state string, r *httpproxy.RequestInstance) (oauth2proxy.TokenInstanceInterface, error)
}

// AuthenticatorInstance is a struct that implements AuthenticatorInstanceInterface.
type AuthenticatorInstance struct {
	FieldAuthenticator *spotifyauth.Authenticator
}

// AuthURL is a proxy for spotifyauth.AuthURL().
func (a *AuthenticatorInstance) AuthURL(state string) string {
	return a.FieldAuthenticator.AuthURL(state)
}

// Token is a proxy for spotifyauth.Authenticator.Token().
func (a *AuthenticatorInstance) Token(ctx *contextproxy.ContextInstance, state string, r *httpproxy.RequestInstance) (oauth2proxy.TokenInstanceInterface, error) {
	token, err := a.FieldAuthenticator.Token(ctx.FieldContext, state, &r.FieldRequest)
	if err != nil {
		return nil, err
	}
	return &oauth2proxy.TokenInstance{FieldToken: *token}, nil
}
