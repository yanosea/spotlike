package oauth2proxy

import (
	"golang.org/x/oauth2"
)

// Oauth2 is an interface for oauth2.
// type Oauth2 interface {
type Oauth2 interface {
	NewToken(tokentype string, refreshToken string) *TokenInstance
}

// Oauth2Proxy is a struct that implements Oauth2.
type Oauth2Proxy struct{}

// New is a constructor for Oauth2Proxy.
func New() Oauth2 {
	return &Oauth2Proxy{}
}

// NewToken is a constructor for TokenInstance.
func (*Oauth2Proxy) NewToken(tokentype string, refreshToken string) *TokenInstance {
	return &TokenInstance{
		FieldToken: &oauth2.Token{
			TokenType:    tokentype,
			RefreshToken: refreshToken,
		},
	}
}
