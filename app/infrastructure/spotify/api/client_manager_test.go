package api

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/yanosea/spotlike/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewClientManager(t *testing.T) {
	origGcm := gcm
	spotify := proxy.NewSpotify()
	http := proxy.NewHttp()
	randstr := proxy.NewRandstr()
	url := proxy.NewUrl()

	type args struct {
		spotify proxy.Spotify
		http    proxy.Http
		randstr proxy.Randstr
		url     proxy.Url
	}
	tests := []struct {
		name    string
		args    args
		want    ClientManager
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (gcm is nil)",
			args: args{
				spotify: spotify,
				http:    http,
				randstr: randstr,
				url:     url,
			},
			want: &clientManager{
				spotify: spotify,
				http:    http,
				randstr: randstr,
				url:     url,
				mutex:   &sync.RWMutex{},
			},
			setup: func() {
				gcm = nil
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name: "positive testing (gcm is not nil)",
			args: args{
				spotify: spotify,
			},
			want: &clientManager{
				spotify: spotify,
				http:    http,
				randstr: randstr,
				url:     url,
				mutex:   &sync.RWMutex{},
			},
			setup: func() {
				gcm = &clientManager{
					spotify: spotify,
					http:    http,
					randstr: randstr,
					url:     url,
					mutex:   &sync.RWMutex{},
				}
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if got := NewClientManager(tt.args.spotify, tt.args.http, tt.args.randstr, tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetClientManager(t *testing.T) {
	origGcm := gcm
	spotify := proxy.NewSpotify()
	http := proxy.NewHttp()
	randstr := proxy.NewRandstr()
	url := proxy.NewUrl()

	tests := []struct {
		name    string
		want    ClientManager
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (gcm is nil)",
			want: nil,
			setup: func() {
				gcm = nil
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name: "positive testing (gcm is not nil)",
			want: &clientManager{
				spotify: spotify,
				http:    http,
				randstr: randstr,
				url:     url,
				mutex:   &sync.RWMutex{},
			},
			setup: func() {
				gcm = &clientManager{
					spotify: spotify,
					http:    http,
					randstr: randstr,
					url:     url,
					mutex:   &sync.RWMutex{},
				}
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if got := GetClientManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClientManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClientManager(t *testing.T) {
	origGcm := gcm
	spotify := proxy.NewSpotify()
	http := proxy.NewHttp()
	randstr := proxy.NewRandstr()
	url := proxy.NewUrl()

	tests := []struct {
		name    string
		want    ClientManager
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (gcm is nil)",
			want: nil,
			setup: func() {
				gcm = nil
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name: "positive testing (gcm is not nil)",
			want: &clientManager{
				spotify: spotify,
				http:    http,
				randstr: randstr,
				url:     url,
				mutex:   &sync.RWMutex{},
			},
			setup: func() {
				gcm = &clientManager{
					spotify: spotify,
					http:    http,
					randstr: randstr,
					url:     url,
					mutex:   &sync.RWMutex{},
				}
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
	}
	for _, tt := range tests {
		if tt.setup != nil {
			tt.setup()
		}
		defer func() {
			if tt.cleanup != nil {
				tt.cleanup()
			}
		}()
		t.Run(tt.name, func(t *testing.T) {
			if got := getClientManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClientManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResetClientManager(t *testing.T) {
	origGcm := gcm
	spotify := proxy.NewSpotify()
	http := proxy.NewHttp()
	randstr := proxy.NewRandstr()
	url := proxy.NewUrl()

	tests := []struct {
		name    string
		isNil   bool
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name:    "positive testing (gcm is nil)",
			isNil:   true,
			wantErr: false,
			setup: func(_ *gomock.Controller) {
				gcm = nil
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name:    "positive testing (gcm is not nil)",
			isNil:   false,
			wantErr: false,
			setup: func(_ *gomock.Controller) {
				gcm = &clientManager{
					spotify: spotify,
					http:    http,
					randstr: randstr,
					url:     url,
					mutex:   &sync.RWMutex{},
				}
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name:    "negative testing (gcm.CloseAllConnections() failed)",
			isNil:   false,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockGcm := NewMockClientManager(mockCtrl)
				mockGcm.EXPECT().CloseClient().Return(errors.New("failed to close client"))
				gcm = mockGcm
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if err := ResetClientManager(); (err != nil) != tt.wantErr {
				t.Errorf("ResetClientManager() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if gcm == nil {
					t.Error("ResetConnectionManager() should not set gcm to nil when error occurs")
				}
			} else {
				if gcm != nil {
					t.Error("ResetConnectionManager() should set gcm to nil when no error occurs")
				}
			}
		})
	}
}

func Test_clientManager_IsClientInitialized(t *testing.T) {
	type fields struct {
		spotify proxy.Spotify
		client  Client
		http    proxy.Http
		randstr proxy.Randstr
		url     proxy.Url
		mutex   *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "positive testing (client is initialized)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  &client{},
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			want: true,
		},
		{
			name: "positive testing (client is not initialized)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientManager{
				spotify: tt.fields.spotify,
				client:  tt.fields.client,
				http:    tt.fields.http,
				randstr: tt.fields.randstr,
				url:     tt.fields.url,
				mutex:   tt.fields.mutex,
			}
			if got := cm.IsClientInitialized(); got != tt.want {
				t.Errorf("clientManager.IsClientInitialized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clientManager_CloseClient(t *testing.T) {
	type fields struct {
		spotify proxy.Spotify
		client  Client
		http    proxy.Http
		randstr proxy.Randstr
		url     proxy.Url
		mutex   *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing (client is nil)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			wantErr: false,
			setup:   nil,
		},
		{
			name: "positive testing (client is not nil)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockClient := NewMockClient(mockCtrl)
				mockClient.EXPECT().Close().Return(nil)
				tt.client = mockClient
			},
		},
		{
			name: "negative testing (client.Close() failed)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockClient := NewMockClient(mockCtrl)
				mockClient.EXPECT().Close().Return(errors.New("failed to close client"))
				tt.client = mockClient
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			cm := &clientManager{
				spotify: tt.fields.spotify,
				client:  tt.fields.client,
				http:    tt.fields.http,
				randstr: tt.fields.randstr,
				url:     tt.fields.url,
				mutex:   tt.fields.mutex,
			}
			if err := cm.CloseClient(); (err != nil) != tt.wantErr {
				t.Errorf("clientManager.CloseClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clientManager_GetClient(t *testing.T) {
	type fields struct {
		spotify proxy.Spotify
		client  Client
		http    proxy.Http
		randstr proxy.Randstr
		url     proxy.Url
		mutex   *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    Client
		wantErr bool
	}{
		{
			name: "positive testing (client is initialized)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  &client{},
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			want:    &client{},
			wantErr: false,
		},
		{
			name: "negative testing (client is not initialized)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientManager{
				spotify: tt.fields.spotify,
				client:  tt.fields.client,
				http:    tt.fields.http,
				randstr: tt.fields.randstr,
				url:     tt.fields.url,
				mutex:   tt.fields.mutex,
			}
			got, err := cm.GetClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("clientManager.GetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clientManager.GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clientManager_InitializeClient(t *testing.T) {
	type fields struct {
		spotify proxy.Spotify
		client  Client
		http    proxy.Http
		randstr proxy.Randstr
		url     proxy.Url
		mutex   *sync.RWMutex
	}
	type args struct {
		ctx    context.Context
		config *ClientConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive testing (client is nil)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  nil,
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			args: args{
				ctx:    context.Background(),
				config: &ClientConfig{},
			},
			wantErr: false,
		},
		{
			name: "negative testing (client is already initialized)",
			fields: fields{
				spotify: proxy.NewSpotify(),
				client:  &client{},
				http:    proxy.NewHttp(),
				randstr: proxy.NewRandstr(),
				url:     proxy.NewUrl(),
				mutex:   &sync.RWMutex{},
			},
			args: args{
				ctx:    context.Background(),
				config: &ClientConfig{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &clientManager{
				spotify: tt.fields.spotify,
				client:  tt.fields.client,
				http:    tt.fields.http,
				randstr: tt.fields.randstr,
				url:     tt.fields.url,
				mutex:   tt.fields.mutex,
			}
			if err := cm.InitializeClient(tt.args.ctx, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("clientManager.InitializeClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
