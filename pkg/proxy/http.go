package proxy

import (
	"context"
	"net/http"
)

// Http is an interface that provides a proxy of the methods of http.
type Http interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	NewServer(addr string) Server
	NotFound(w http.ResponseWriter, r *http.Request)
}

// ResponseWriter is an interface that provides a proxy of the methods of http.ResponseWriter.
type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

// httpProxy is a proxy struct that implements the Http interface.
type httpProxy struct{}

// NewHttp returns a new instance of the Http interface.
func NewHttp() Http {
	return &httpProxy{}
}

// NewResponseWriter returns a new instance of the ResponseWriter interface.
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriterProxy{w: w}
}

// responseWriterProxy is a proxy struct that implements the ResponseWriter interface.
type responseWriterProxy struct {
	w http.ResponseWriter
}

// Header is a proxy method that calls the Header method of the http.ResponseWriter.
func (r *responseWriterProxy) Header() http.Header {
	return r.w.Header()
}

// Write is a proxy method that calls the Write method of the http.ResponseWriter.
func (r *responseWriterProxy) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

// WriteHeader is a proxy method that calls the WriteHeader method of the http.ResponseWriter.
func (r *responseWriterProxy) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
}

// HandleFunc is a proxy method that calls the HandleFunc method of the http.
func (h *httpProxy) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}

// NewServer is a proxy method that returns the http.Server.
func (*httpProxy) NewServer(addr string) Server {
	return &serverProxy{server: &http.Server{Addr: addr}}
}

// NotFound is a proxy method that calls the NotFound method of the http.
func (*httpProxy) NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

// Server is an interface that provides a proxy of the methods of http.Server.
type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// serverProxy is a proxy struct that implements the Server interface.
type serverProxy struct {
	server *http.Server
}

// ListenAndServe is a proxy method that calls the ListenAndServe method of the http.Server.
func (s *serverProxy) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Shutdown is a proxy method that calls the Shutdown method of the http.Server.
func (s *serverProxy) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
