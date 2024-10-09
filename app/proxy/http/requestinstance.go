package httpproxy

import (
	"net/http"

	"github.com/yanosea/spotlike/app/proxy/context"
)

// RequestInstanceInterface is an interface for http.Request.
type RequestInstanceInterface interface {
	Context() contextproxy.ContextInstance
}

// RequestInstance is a struct that implements RequestInstanceInterface.
type RequestInstance struct {
	FieldRequest http.Request
}

func (r *RequestInstance) Context() *contextproxy.ContextInstance {
	return &contextproxy.ContextInstance{FieldContext: r.FieldRequest.Context()}
}
