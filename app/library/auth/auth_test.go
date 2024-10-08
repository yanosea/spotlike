package auth

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/yanosea/spotlike/util"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func TestIsEnvsSet(t *testing.T) {
	tests := []struct {
		name                string
		spotifyID           string
		spotifySecret       string
		spotifyRedirectUri  string
		spotifyRefreshToken string
		want                bool
	}{
		{
			name:                "positive testing (all set)",
			spotifyID:           "test_id",
			spotifySecret:       "test_secret",
			spotifyRedirectUri:  "test_redirect_uri",
			spotifyRefreshToken: "test_refresh_token",
			want:                true,
		}, {
			name:                "positive testing (SPOTIFY_ID not set)",
			spotifyID:           "",
			spotifySecret:       "test_secret",
			spotifyRedirectUri:  "test_redirect_uri",
			spotifyRefreshToken: "test_refresh_token",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_SECRET not set)",
			spotifyID:           "test_id",
			spotifySecret:       "",
			spotifyRedirectUri:  "test_redirect_uri",
			spotifyRefreshToken: "test_refresh_token",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_REDIRECT_URI not set)",
			spotifyID:           "test_id",
			spotifySecret:       "test_secret",
			spotifyRedirectUri:  "",
			spotifyRefreshToken: "test_refresh_token",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_REFRESH_TOKEN not set)",
			spotifyID:           "test_id",
			spotifySecret:       "test_secret",
			spotifyRedirectUri:  "test_redirect_uri",
			spotifyRefreshToken: "",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_ID and SPOTIFY_SECRET not set)",
			spotifyID:           "",
			spotifySecret:       "",
			spotifyRedirectUri:  "test_redirect_uri",
			spotifyRefreshToken: "test_refresh_token",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_ID and SPOTIFY_REDIRECT_URI not set)",
			spotifyID:           "",
			spotifySecret:       "test_secret",
			spotifyRedirectUri:  "",
			spotifyRefreshToken: "test_refresh_token",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_ID and SPOTIFY_REFRESH_TOKEN not set)",
			spotifyID:           "",
			spotifySecret:       "test_secret",
			spotifyRedirectUri:  "test_redirect_uri",
			spotifyRefreshToken: "",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_SECRET and SPOTIFY_REDIRECT_URI not set)",
			spotifyID:           "test_id",
			spotifySecret:       "",
			spotifyRedirectUri:  "",
			spotifyRefreshToken: "test_refresh_token",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_SECRET and SPOTIFY_REFRESH_TOKEN not set)",
			spotifyID:           "test_id",
			spotifySecret:       "",
			spotifyRedirectUri:  "test_redirect_uri",
			spotifyRefreshToken: "",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_REDIRECT_URI and SPOTIFY_REFRESH_TOKEN not set)",
			spotifyID:           "test_id",
			spotifySecret:       "test_secret",
			spotifyRedirectUri:  "",
			spotifyRefreshToken: "",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_ID, SPOTIFY_SECRET, SPOTIFY_REDIRECT_URI not set)",
			spotifyID:           "",
			spotifySecret:       "",
			spotifyRedirectUri:  "",
			spotifyRefreshToken: "test_refresh_token",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_ID, SPOTIFY_SECRET, SPOTIFY_REFRESH_TOKEN not set)",
			spotifyID:           "",
			spotifySecret:       "",
			spotifyRedirectUri:  "test_redirect_uri",
			spotifyRefreshToken: "",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_ID, SPOTIFY_REDIRECT_URI, SPOTIFY_REFRESH_TOKEN not set)",
			spotifyID:           "",
			spotifySecret:       "test_secret",
			spotifyRedirectUri:  "",
			spotifyRefreshToken: "",
			want:                false,
		}, {
			name:                "positive testing (SPOTIFY_SECRET, SPOTIFY_REDIRECT_URI, SPOTIFY_REFRESH_TOKEN not set)",
			spotifyID:           "test_id",
			spotifySecret:       "",
			spotifyRedirectUri:  "",
			spotifyRefreshToken: "",
			want:                false,
		}, {
			name:                "positive testing (all empty)",
			spotifyID:           "",
			spotifySecret:       "",
			spotifyRedirectUri:  "",
			spotifyRefreshToken: "",
			want:                false,
		},
	}

	for _, tt := range tests {
		os.Setenv(util.AUTH_ENV_SPOTIFY_ID, tt.spotifyID)
		os.Setenv(util.AUTH_ENV_SPOTIFY_SECRET, tt.spotifySecret)
		os.Setenv(util.AUTH_ENV_SPOTIFY_REDIRECT_URI, tt.spotifyRedirectUri)
		os.Setenv(util.AUTH_ENV_SPOTIFY_REFRESH_TOKEN, tt.spotifyRefreshToken)
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEnvsSet(); got != tt.want {
				t.Errorf("IsEnvsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAuthInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAuthInfo()
		})
	}
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name    string
		want    *spotify.Client
		wantO   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &bytes.Buffer{}
			got, err := Authenticate(o)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authenticate() = %v, want %v", got, tt.want)
			}
			if gotO := o.String(); gotO != tt.wantO {
				t.Errorf("Authenticate() = %v, want %v", gotO, tt.wantO)
			}
		})
	}
}

func Test_refresh(t *testing.T) {
	type args struct {
		authenticator *spotifyauth.Authenticator
		refreshToken  string
	}
	tests := []struct {
		name    string
		args    args
		want    *spotify.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := refresh(tt.args.authenticator, tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("refresh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("refresh() = %v, want %v", got, tt.want)
			}
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
			}, want: "8080",
			wantErr: false,
		}, {
			name: "negative testing (the uri doesn't have port number)",
			args: args{
				uri: "http://localhost/callback",
			}, wantErr: true,
			wantMsg: auth_error_message_invalid_uri,
		}, {
			name: "negative testing (space)",
			args: args{
				uri: " ",
			}, wantErr: true,
			wantMsg: auth_error_message_invalid_uri,
		}, {
			name: "negative testing (blank)",
			args: args{
				uri: "",
			}, wantErr: true,
			wantMsg: auth_error_message_invalid_uri,
		}, {
			name: "negative testing (invalid uri)",
			args: args{
				uri: "http://localhost:8080/invalid%2",
			}, wantErr: true,
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
