package oauth2proxy

import ()

// Oauth2 is an interface for oauth2.
// type Oauth2 interface {
type Oauth2 interface{}

// Oauth2Proxy is a struct that implements Oauth2.
type Oauth2Proxy struct{}

// New is a constructor for Oauth2Proxy.
func New() Oauth2 {
	return &Oauth2Proxy{}
}
