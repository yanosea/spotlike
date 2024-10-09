package httpproxy

import (
	"net/http"
)

// ResponseWriterInstanceInterface is an interface for http.ResponseWriter.
type ResponseWriterInstanceInterface interface {
	http.ResponseWriter
	Header() http.Header
}

// ResponseWriterInstance is a struct that implements ResponseWriterInstanceInterface.
type ResponseWriterInstance struct {
	FieldResponseWriter http.ResponseWriter
}

// Header is a proxy for http.ResponseWriter.Header().
func (r *ResponseWriterInstance) Header() http.Header {
	return r.FieldResponseWriter.Header()
}

// Write is a proxy for http.ResponseWriter.Write().
func (r *ResponseWriterInstance) Write(b []byte) (int, error) {
	return r.FieldResponseWriter.Write(b)
}

// WriteHeader is a proxy for http.ResponseWriter.WriteHeader().
func (r *ResponseWriterInstance) WriteHeader(statusCode int) {
	r.FieldResponseWriter.WriteHeader(statusCode)
}
