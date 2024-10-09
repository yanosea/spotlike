package oauth2proxy

import (
	"golang.org/x/oauth2"
)

// TokenInstanceInterface is an interface for TokenInstance.
type TokenInstanceInterface interface{}

// TokenInstance is a struct that implements TokenInstanceInterface.
type TokenInstance struct {
	FieldToken *oauth2.Token
}

// NewToken is a constructor for TokenInstance.
func NewToken(tokentype string, refreshToken string) TokenInstanceInterface {
	return &TokenInstance{
		FieldToken: &oauth2.Token{
			TokenType:    tokentype,
			RefreshToken: refreshToken,
		},
	}
}
