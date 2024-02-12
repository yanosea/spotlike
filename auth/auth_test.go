package auth

import (
	"crypto/rand"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/zmb3/spotify/v2"
)

func TestGetClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *spotify.Client
		wantErr bool
	}{
		{
			name: "positive testing",
			want: &spotify.Client{},
		}, {
			name:    "negative testing (authenticate error)",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		sa := &SpotlikeAuthenticator{}
		if tt.wantErr {
			fakeSpotlikeAuthenticator := &FakeSpotlikeAuthenticator{
				FakeGetClient: func() (*spotify.Client, error) {
					return nil, errors.New("fake GetClient error")
				},
			}
			sa = &SpotlikeAuthenticator{c: fakeSpotlikeAuthenticator}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := sa.GetClient()
			if err == nil && reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GetClient() = %v, want = %v", got, tt.want)
				return
			}
			if err != nil && tt.wantErr {
			}
		})
	}
}

type FakeSpotlikeAuthenticator struct {
	Authenticator
	FakeGetClient func() (*spotify.Client, error)
}

func TestSetAuthInfo(t *testing.T) {
	tests := []struct {
		name string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name    string
		want    *spotify.Client
		wantErr bool
	}{}
	for _, tt := range tests {
		sa := &SpotlikeAuthenticator{}
		t.Run(tt.name, func(t *testing.T) {
			got, err := sa.authenticate()
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

func TestCompleteAuthenticate(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		sa := &SpotlikeAuthenticator{}
		t.Run(tt.name, func(t *testing.T) {
			sa.completeAuthenticate(tt.args.w, tt.args.r)
		})
	}
}
func TestGetPort(t *testing.T) {
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

func TestGenerateRandomString(t *testing.T) {
	type args struct {
		length int
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantReadErr bool
		wantMsg     string
	}{
		{
			name: "positive testing",
			args: args{
				length: 16,
			},
			wantErr:     false,
			wantReadErr: false,
		}, {
			name: "positive testing (zero)",
			args: args{
				length: 0,
			},
			wantErr:     false,
			wantReadErr: false,
		}, {
			name: "negative testing (negative length)",
			args: args{
				length: -1,
			},
			wantErr:     true,
			wantReadErr: false,
			wantMsg:     auth_error_message_invalid_length_for_random_string,
		}, {
			name: "negative testing (failed to generate a random number)",
			args: args{
				length: 16,
			},
			wantErr:     true,
			wantReadErr: true,
			wantMsg:     "fake Read error",
		},
	}
	for _, tt := range tests {
		if tt.wantReadErr {
			oldRandReader := rand.Reader
			rand.Reader = &fakeRandReader{}
			defer func() { rand.Reader = oldRandReader }()
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateRandomString(tt.args.length)
			if err != nil && !tt.wantErr && tt.args.length > 0 && len(got) != tt.args.length {
				t.Errorf("GenerateRandomString() = length %v, want length %v", len(got), tt.args.length)
				return
			}
			if err != nil && tt.wantErr {
				msg := err.Error()
				if !strings.Contains(msg, tt.wantMsg) {
					t.Errorf("GenerateRandomString() unexpected error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if err != nil && !tt.wantErr {
				t.Errorf("GenerateRandomString() unexpected error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type fakeRandReader struct{}

func (f *fakeRandReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("fake Read error")
}
