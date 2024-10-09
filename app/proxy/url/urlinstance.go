package urlproxy

import (
	"net/url"
)

// UrlInstanceInterface is an interface for url.URL.
type UrlInstanceInterface interface {
	Port() string
}

// UrlInstance is a struct that implements Url.
type UrlInstance struct {
	FieldUrl *url.URL
}

// Port is a proxy for url.URL.Port().
func (u *UrlInstance) Port() string {
	return u.FieldUrl.Port()
}
