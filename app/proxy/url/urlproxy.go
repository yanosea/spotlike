package urlproxy

import (
	"net/url"
)

// Url is an interface for Url.
type Url interface {
	Parse(rawURL string) (UrlInstanceInterface, error)
}

// UrlProxy is a struct that implements Url.
type UrlProxy struct{}

// New is a constructor for UrlProxy.
func New() Url {
	return &UrlProxy{}
}

// Parse is a proxy for url.Parse.
func (*UrlProxy) Parse(rawURL string) (UrlInstanceInterface, error) {
	url, _ := url.Parse(rawURL)
	return &UrlInstance{FieldUrl: url}, nil
}
