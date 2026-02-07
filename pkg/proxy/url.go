package proxy

import (
	"net/url"
)

// Url is an interface that provides a proxy of the methods of url.
type Url interface {
	Parse(rawurl string) (*url.URL, error)
}

// urlProxy is a proxy struct that implements the Url interface.
type urlProxy struct{}

// NewUrl returns a new instance of the Url interface.
func NewUrl() Url {
	return &urlProxy{}
}

// Parse is a proxy method that calls the Parse method of the url.
func (*urlProxy) Parse(rawurl string) (*url.URL, error) {
	return url.Parse(rawurl)
}
