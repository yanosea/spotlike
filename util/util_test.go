package util

import (
	"testing"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

func TestGetPortString(t *testing.T) {
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
			want:    ":8080",
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
				uri: "test:8080",
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
			got, err := GetPortStringFromUri(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPortString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPortString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive testing",
			args: args{
				length: 1,
			},
			wantErr: false,
		},
		{
			name: "positive testing (zero)",
			args: args{
				length: 0,
			},
			wantErr: false,
		},
		{
			name: "negative testing (negative length)",
			args: args{
				length: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomString(tt.args.length)
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
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive testing",
			args: args{
				artists: []spotify.SimpleArtist{{Name: "TestArtist1"}, {Name: "TestArtist2"}},
			},
			want:    "TestArtist1, TestArtist2",
			wantErr: false,
		},
		{
			name: "positive testing (solo)",
			args: args{
				artists: []spotify.SimpleArtist{{Name: "TestArtist1"}},
			},
			want:    "TestArtist1",
			wantErr: false,
		},
		{
			name: "positive testing (blank)",
			args: args{
				artists: []spotify.SimpleArtist{},
			},
			wantErr: true,
		},
		{
			name: "negative testing (nil)",
			args: args{
				artists: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CombineArtistNames(tt.args.artists)
			if (err != nil) != tt.wantErr {
				t.Errorf("CombineArtistNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CombineArtistNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
