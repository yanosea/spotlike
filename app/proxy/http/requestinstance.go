package httpproxy

import (
	"net/http"

	"github.com/yanosea/spotlike/app/proxy/context"
)

// RequestInstanceInterface is an interface for http.Request.
type RequestInstanceInterface interface {
	Context() contextproxy.ContextInstance
	FormValue(key string) string
}

// RequestInstance is a struct that implements RequestInstanceInterface.
type RequestInstance struct {
	FieldRequest http.Request
}

// Context is a proxy for http.Request.Context().
func (r *RequestInstance) Context() *contextproxy.ContextInstance {
	return &contextproxy.ContextInstance{FieldContext: r.FieldRequest.Context()}
}

// FormValue is a proxy for http.Request.FormValue().
func (r *RequestInstance) FormValue(key string) string {
	return r.FieldRequest.FormValue(key)
}
