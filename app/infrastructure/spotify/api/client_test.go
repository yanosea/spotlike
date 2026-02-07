package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/yanosea/spotlike/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func Test_client_Close(t *testing.T) {
	type fields struct {
		spotify proxy.Spotify
		client  proxy.Client
		http    proxy.Http
		url     proxy.Url
		config  *ClientConfig
		context context.Context
		mutex   *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (client is not nil)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				url:     proxy.NewUrl(),
				config:  &ClientConfig{},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockClient := proxy.NewMockClient(mockCtrl)
				tt.client = mockClient
			},
		},
		{
			name: "positive testing (client is nil)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				url:     proxy.NewUrl(),
				config:  &ClientConfig{},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			wantErr: false,
			setup:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			c := &client{
				spotify: tt.fields.spotify,
				client:  tt.fields.client,
				http:    tt.fields.http,
				url:     tt.fields.url,
				config:  tt.fields.config,
				context: tt.fields.context,
				mutex:   tt.fields.mutex,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("client.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Open(t *testing.T) {
	clientConfig := &ClientConfig{
		SpotifyID:           "test-id",
		SpotifySecret:       "test-secret",
		SpotifyRedirectUri:  "http://localhost:8080/callback",
		SpotifyRefreshToken: "test-refresh-token",
	}

	type fields struct {
		spotify proxy.Spotify
		client  proxy.Client
		http    proxy.Http
		url     proxy.Url
		config  *ClientConfig
		context context.Context
		mutex   *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		setup  func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (client is already initialized)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    proxy.NewHttp(),
				url:     proxy.NewUrl(),
				config:  clientConfig,
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockClient := proxy.NewMockClient(mockCtrl)
				tt.client = mockClient
			},
		},
		{
			name: "positive testing (client is not initialized)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				url:     proxy.NewUrl(),
				config:  clientConfig,
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			c := &client{
				spotify: tt.fields.spotify,
				client:  tt.fields.client,
				http:    tt.fields.http,
				url:     tt.fields.url,
				config:  tt.fields.config,
				context: tt.fields.context,
				mutex:   tt.fields.mutex,
			}
			initialClient := c.client
			got := c.Open()
			if got == nil {
				t.Errorf("client.Open() returned nil, expected non-nil client")
			}
			if initialClient != nil {
				if got != initialClient {
					t.Errorf("client.Open() = %v, expected same instance as initialClient %v", got, initialClient)
				}
			} else {
				if c.client == nil {
					t.Errorf("client.client field is nil after c.Open(), expected to be initialized")
				}
				if got != c.client {
					t.Errorf("client.Open() = %v, expected same as c.client %v", got, c.client)
				}
			}
		})
	}
}

func Test_client_UpdateConfig(t *testing.T) {
	type fields struct {
		spotify proxy.Spotify
		client  proxy.Client
		http    proxy.Http
		url     proxy.Url
		config  *ClientConfig
		context context.Context
		mutex   *sync.RWMutex
	}
	type args struct {
		config *ClientConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive testing",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				url:     proxy.NewUrl(),
				config: &ClientConfig{
					SpotifyID:           "old-id",
					SpotifySecret:       "old-secret",
					SpotifyRedirectUri:  "http://localhost:8080/old-callback",
					SpotifyRefreshToken: "old-refresh-token",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			args: args{
				config: &ClientConfig{
					SpotifyID:           "new-id",
					SpotifySecret:       "new-secret",
					SpotifyRedirectUri:  "http://localhost:8080/new-callback",
					SpotifyRefreshToken: "new-refresh-token",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				spotify: tt.fields.spotify,
				client:  tt.fields.client,
				http:    tt.fields.http,
				url:     tt.fields.url,
				config:  tt.fields.config,
				context: tt.fields.context,
				mutex:   tt.fields.mutex,
			}
			origConfig := c.config
			c.UpdateConfig(tt.args.config)
			if reflect.DeepEqual(origConfig, c.config) {
				t.Errorf("client.UpdateConfig() updated config to %v, expected %v", c.config, tt.args.config)
			}
		})
	}
}

func Test_client_Auth(t *testing.T) {
	type fields struct {
		spotify proxy.Spotify
		client  proxy.Client
		http    proxy.Http
		randstr proxy.Randstr
		url     proxy.Url
		config  *ClientConfig
		context context.Context
		mutex   *sync.RWMutex
	}
	type args struct {
		authUrlChan chan<- string
		version     string
	}
	tests := []struct {
		name        string
		fields      fields
		setup       func(mockCtrl *gomock.Controller, tt *fields)
		args        args
		wantClient  bool
		wantAuthUrl bool
		wantErr     bool
	}{
		{
			name: "positive testing (client is already initialized)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				config:  &ClientConfig{},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				tt.client = mockSpotifyClient
			},
			wantClient:  true,
			wantAuthUrl: false,
			wantErr:     false,
		},
		{
			name: "positive testing (client is not initialized)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockAuthenticator.EXPECT().Token(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&oauth2.Token{
						RefreshToken: "test-refresh-token",
					},
					nil,
				)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockResponseWriter := proxy.NewMockResponseWriter(mockCtrl)
				mockResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()
				mockResponseWriter.EXPECT().WriteHeader(http.StatusOK)
				mockResponseWriter.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).Do(
					func(pattern string, handler http.HandlerFunc) {
						go func() {
							req, _ := http.NewRequest("GET", "/callback", nil)
							handler(mockResponseWriter, req)
						}()
					},
				)
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  true,
			wantAuthUrl: true,
			wantErr:     false,
		},
		{
			name: "negative testing (c.url.Parse() failed)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "invalid:uri",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse("invalid:uri").Return(nil, errors.New("failed to parse"))
				tt.url = mockUrl
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (url.Port() == \"\")",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse("http://localhost/callback").Return(&url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "/callback",
				}, nil)
				tt.url = mockUrl
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (auenticatior.Token() failed)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockAuthenticator.EXPECT().Token(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get token"))
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockResponseWriter := proxy.NewMockResponseWriter(mockCtrl)
				mockResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()
				mockResponseWriter.EXPECT().WriteHeader(http.StatusForbidden)
				mockResponseWriter.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).Do(
					func(pattern string, handler http.HandlerFunc) {
						go func() {
							req, _ := http.NewRequest("GET", "/callback", nil)
							handler(mockResponseWriter, req)
						}()
					},
				)
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (tok == nil || tok.RefreshToken == \"\")",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockAuthenticator.EXPECT().Token(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&oauth2.Token{
						TokenType:    "Bearer",
						AccessToken:  "test-access-token",
						RefreshToken: "",
					},
					nil,
				)
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockResponseWriter := proxy.NewMockResponseWriter(mockCtrl)
				mockResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()
				mockResponseWriter.EXPECT().WriteHeader(http.StatusForbidden)
				mockResponseWriter.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).Do(
					func(pattern string, handler http.HandlerFunc) {
						go func() {
							req, _ := http.NewRequest("GET", "/callback", nil)
							handler(mockResponseWriter, req)
						}()
					},
				)
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (w.Write() failed)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().Token(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&oauth2.Token{
						RefreshToken: "test-refresh-token",
					}, nil,
				)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockResponseWriter := proxy.NewMockResponseWriter(mockCtrl)
				mockResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()
				mockResponseWriter.EXPECT().WriteHeader(http.StatusOK)
				mockResponseWriter.EXPECT().Write(gomock.Any()).Return(0, errors.New("failed to write"))
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).Do(
					func(pattern string, handler http.HandlerFunc) {
						go func() {
							req, _ := http.NewRequest("GET", "/callback", nil)
							handler(mockResponseWriter, req)
						}()
					},
				)
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (server.ListenAndServe() failed)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().Return(errors.New("failed to listen and serve"))
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any())
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (context timeout)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				mutex: &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (server.Shutdown() failed)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockAuthenticator.EXPECT().Token(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&oauth2.Token{
						RefreshToken: "test-refresh-token",
					},
					nil,
				)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockResponseWriter := proxy.NewMockResponseWriter(mockCtrl)
				mockResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()
				mockResponseWriter.EXPECT().WriteHeader(http.StatusOK)
				mockResponseWriter.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("failed to shutdown"))
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).Do(
					func(pattern string, handler http.HandlerFunc) {
						go func() {
							req, _ := http.NewRequest("GET", "/callback", nil)
							handler(mockResponseWriter, req)
						}()
					},
				)
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (server shutdown timeout)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockAuthenticator.EXPECT().Token(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&oauth2.Token{
						RefreshToken: "test-refresh-token",
					},
					nil,
				)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockResponseWriter := proxy.NewMockResponseWriter(mockCtrl)
				mockResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()
				mockResponseWriter.EXPECT().WriteHeader(http.StatusOK)
				mockResponseWriter.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()
				mockServer := proxy.NewMockServer(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).Do(
					func(pattern string, handler http.HandlerFunc) {
						go func() {
							req, _ := http.NewRequest("GET", "/callback", nil)
							handler(mockResponseWriter, req)
						}()
					},
				)
				mockServer.EXPECT().ListenAndServe().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
					<-ctx.Done()
					return ctx.Err()
				})
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (nil client received)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockAuthenticator.EXPECT().Token(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&oauth2.Token{
						RefreshToken: "test-refresh-token",
					},
					nil,
				)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(nil)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockResponseWriter := proxy.NewMockResponseWriter(mockCtrl)
				mockResponseWriter.EXPECT().Header().Return(http.Header{}).AnyTimes()
				mockResponseWriter.EXPECT().WriteHeader(http.StatusOK)
				mockResponseWriter.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).Do(
					func(pattern string, handler http.HandlerFunc) {
						go func() {
							req, _ := http.NewRequest("GET", "/callback", nil)
							handler(mockResponseWriter, req)
						}()
					},
				)
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (client not received)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().DoAndReturn(func() error {
					time.Sleep(10 * time.Millisecond)
					go func() {
						time.Sleep(100 * time.Millisecond)
						select {
						case dc <- struct{}{}:
						default:
						}
					}()
					return nil
				})
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).Do(
					func(pattern string, handler http.HandlerFunc) {
					},
				)
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (context timeout with cleanup error)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				mutex: &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().AnyTimes().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("timeout to shutdown"))
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any()).AnyTimes()
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan<- string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
		{
			name: "negative testing (error chan with cleanup error)",
			fields: fields{
				spotify: nil,
				client:  nil,
				http:    nil,
				randstr: proxy.NewRandstr(),
				url:     nil,
				config: &ClientConfig{
					SpotifyRedirectUri: "http://localhost:8080/callback",
				},
				context: context.Background(),
				mutex:   &sync.RWMutex{},
			},
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				tt.spotify = mockSpotify
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(&url.URL{Host: "localhost:8080"}, nil)
				tt.url = mockUrl
				mockServer := proxy.NewMockServer(mockCtrl)
				mockServer.EXPECT().ListenAndServe().Return(errors.New("failed to listen and serve"))
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(errors.New("failed to shutdown"))
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().NewServer(":8080").Return(mockServer)
				mockHttp.EXPECT().HandleFunc("/callback", gomock.Any())
				tt.http = mockHttp
			},
			args: args{
				authUrlChan: make(chan string, 1),
				version:     "v0.0.0",
			},
			wantClient:  false,
			wantAuthUrl: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			var authChan chan string
			if tt.args.authUrlChan != nil {
				authChan = make(chan string, 1)
				tt.args.authUrlChan = authChan
			}
			c := &client{
				spotify: tt.fields.spotify,
				client:  tt.fields.client,
				http:    tt.fields.http,
				randstr: tt.fields.randstr,
				url:     tt.fields.url,
				config:  tt.fields.config,
				context: tt.fields.context,
				mutex:   tt.fields.mutex,
			}
			got, got1, err := c.Auth(tt.args.authUrlChan, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantClient {
				t.Errorf("client.Auth() got = %v, wantClient %v", got, tt.wantClient)
			}
			if (got1 != "") != tt.wantAuthUrl {
				t.Errorf("client.Auth() got1 = %v, wantAuthUrl %v", got1, tt.wantAuthUrl)
			}
		})
	}
}
