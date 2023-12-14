package api

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/zmb3/spotify/v2"
)

func TestGetClient(t *testing.T) {
	tests := []struct {
		name               string
		want               *spotify.Client
		wantErr            bool
		wantSetAuthInfoErr bool
		wantMsg            string
	}{
		{
			name: "positive testing",
			want: &spotify.Client{},
		},
		{
			name:               "negative testing (setAuthInfo error)",
			wantErr:            true,
			wantSetAuthInfoErr: true,
			wantMsg:            "fake setAuthInfo error",
		},
	}
	for _, tt := range tests {
		if tt.wantSetAuthInfoErr {
			oldSetAuthInfo := setAuthInfo
			setAuthInfo = setAuthInfoFake
			defer func() { setAuthInfo = oldSetAuthInfo }()
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClient()
			if err == nil && reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GetClient() = %v, want = %v", got, tt.want)
				return
			}
			if err != nil && tt.wantErr {
				msg := err.Error()
				if !strings.Contains(msg, tt.wantMsg) {
					t.Errorf("GetClient() unexpected error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if err != nil && !tt.wantErr {
				t.Errorf("GetClient() unexpected error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_setAuthInfo(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setAuthInfo(); (err != nil) != tt.wantErr {
				t.Errorf("setAuthInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authenticate(t *testing.T) {
	tests := []struct {
		name    string
		want    *spotify.Client
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authenticate()
			if (err != nil) != tt.wantErr {
				t.Errorf("authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_completeAuthenticate(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completeAuthenticate(tt.args.w, tt.args.r)
		})
	}
}

func setAuthInfoFake() error {
	return errors.New("fake setAuthInfo error")
}
