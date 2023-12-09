package api

import (
	"crypto/rand"
	"errors"
	"testing"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

func TestGetPort(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive testing",
			args: args{
				uri: "http://localhost:8080/callback",
			},
			want:    "8080",
			wantErr: false,
		},
		{
			name: "negative testing (the uri doesn't have port number)",
			args: args{
				uri: "http://localhost/callback",
			},
			wantErr: true,
		},
		{
			name: "negative testing (invalid uri)",
			args: args{
				uri: "http://localhost:8080/invalid%2",
			},
			wantErr: true,
		},
		{
			name: "negative testing (blank)",
			args: args{
				uri: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPortFromUri(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPortFromUri() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPortFromUri() = %v, want %v", got, tt.want)
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
	}{
		{
			name: "positive testing",
			args: args{
				length: 16,
			},
			wantErr:     false,
			wantReadErr: false,
		},
		{
			name: "positive testing (zero)",
			args: args{
				length: 0,
			},
			wantErr:     false,
			wantReadErr: false,
		},
		{
			name: "negative testing (negative length)",
			args: args{
				length: -1,
			},
			wantErr:     true,
			wantReadErr: false,
		},
		{
			name: "negative testing (failed to generate a random number)",
			args: args{
				length: 16,
			},
			wantErr:     true,
			wantReadErr: true,
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
			if (err != nil) != tt.wantErr && tt.wantReadErr {
				t.Errorf("GenerateRandomString() error = %v, wantReadErr %v", err, tt.wantReadErr)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.length > 0 && len(got) != tt.args.length {
				t.Errorf("GenerateRandomString() = length %v, want length %v", len(got), tt.args.length)
			}
		})
	}
}

func TestCombineArtistNames(t *testing.T) {
	type args struct {
		artists []spotify.SimpleArtist
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive testing",
			args: args{
				artists: []spotify.SimpleArtist{{Name: "TestArtist1"}, {Name: "TestArtist2"}},
			},
			want: "TestArtist1, TestArtist2",
		},
		{
			name: "positive testing (solo)",
			args: args{
				artists: []spotify.SimpleArtist{{Name: "TestArtist1"}},
			},
			want: "TestArtist1",
		},
		{
			name: "positive testing (blank)",
			args: args{
				artists: []spotify.SimpleArtist{},
			},
			want: "",
		},
		{
			name: "negative testing (nil)",
			args: args{
				artists: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := combineArtistNames(tt.args.artists)
			if got != tt.want {
				t.Errorf("CombineArtistNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

type fakeRandReader struct{}

func (f *fakeRandReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("fake Read error")
}
