package httpproxy

import (
	"net/http"
)

// Http is an interface for http.
type Http interface {
	Error(w ResponseWriterInstanceInterface, error error, code int)
	HandleFunc(pattern string, handler func(ResponseWriterInstanceInterface, *RequestInstance))
}

// HttpProxy is a struct that implements Http.
type HttpProxy struct{}

// New is a constructor for HttpProxy.
func New() Http {
	return &HttpProxy{}
}

func (*HttpProxy) Error(w ResponseWriterInstanceInterface, error error, code int) {
	http.Error(w, error.Error(), code)
}

// HandleFunc is a proxy for http.HandleFunc.
func (*HttpProxy) HandleFunc(pattern string, handler func(ResponseWriterInstanceInterface, *RequestInstance)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handler(&ResponseWriterInstance{FieldResponseWriter: w}, &RequestInstance{FieldRequest: *r})
	})
}
