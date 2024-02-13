package auth

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func()
		want     *authenticator
	}{
		{
			name: "Positive case",
			setupEnv: func() {
				os.Setenv(auth_env_spotify_id, "test_spotify_id")
				os.Setenv(auth_env_spotify_secret, "test_spotify_secret")
				os.Setenv(auth_env_spotify_redirect_uri, "http://localhost:8080/callback")
				os.Setenv(auth_env_spotify_refresh_token, "")
			},
			want: &authenticator{
				spotifyAuthenticator: spotifyauth.New(
					spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
					spotifyauth.WithScopes(
						spotifyauth.ScopeUserFollowRead,
						spotifyauth.ScopeUserLibraryRead,
						spotifyauth.ScopeUserFollowModify,
						spotifyauth.ScopeUserLibraryModify,
					),
				),
				channel: make(chan *spotify.Client),
				state:   "test_state",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			got := New()
			fmt.Println(tt.want)
			fmt.Println(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setAuthInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setAuthInfo()
		})
	}
}

func Test_authenticator_Authenticate(t *testing.T) {
	type fields struct {
		spotifyAuthenticator *spotifyauth.Authenticator
		channel              chan *spotify.Client
		state                string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *spotify.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authenticator{
				spotifyAuthenticator: tt.fields.spotifyAuthenticator,
				channel:              tt.fields.channel,
				state:                tt.fields.state,
			}
			got, err := a.Authenticate()
			if (err != nil) != tt.wantErr {
				t.Errorf("authenticator.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authenticator.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authenticator_completeAuthenticate(t *testing.T) {
	type fields struct {
		spotifyAuthenticator *spotifyauth.Authenticator
		channel              chan *spotify.Client
		state                string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authenticator{
				spotifyAuthenticator: tt.fields.spotifyAuthenticator,
				channel:              tt.fields.channel,
				state:                tt.fields.state,
			}
			a.completeAuthenticate(tt.args.w, tt.args.r)
		})
	}
}

func Test_getPortFromUri(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		wantMsg string
	}{
		{
			name: "positive testing",
			args: args{
				uri: "http://localhost:8080/callback",
			},
			want:    "8080",
			wantErr: false,
		}, {
			name: "negative testing (the uri doesn't have port number)",
			args: args{
				uri: "http://localhost/callback",
			},
			wantErr: true,
			wantMsg: auth_error_message_invalid_uri,
		}, {
			name: "negative testing (space)",
			args: args{
				uri: " ",
			},
			wantErr: true,
			wantMsg: auth_error_message_invalid_uri,
		}, {
			name: "negative testing (blank)",
			args: args{
				uri: "",
			},
			wantErr: true,
			wantMsg: auth_error_message_invalid_uri,
		}, {
			name: "negative testing (invalid uri)",
			args: args{
				uri: "http://localhost:8080/invalid%2",
			},
			wantErr: true,
			wantMsg: "invalid URL escape",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPortFromUri(tt.args.uri)
			if err != nil && !tt.wantErr && got != tt.want {
				t.Errorf("GetPortFromUri() = %v, want %v", got, tt.want)
				return
			}
			if err != nil && tt.wantErr {
				msg := err.Error()
				if !strings.Contains(msg, tt.wantMsg) {
					t.Errorf("GetPortFromUri() unexpected error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if err != nil && !tt.wantErr {
				t.Errorf("GetPortFromUri() unexpected error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
