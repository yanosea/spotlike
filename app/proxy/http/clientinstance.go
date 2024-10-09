package httpproxy

import (
	"net/http"
)

// ClientInstanceInterface is an interface for http.Client.
type ClientInstanceInterface interface{}

// ClientInstance is a struct that implements ClientInstanceInterface.
type ClientInstance struct {
	FieldClient *http.Client
}
