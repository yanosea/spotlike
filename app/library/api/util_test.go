package api

import (
	"testing"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

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
		}, {
			name: "positive testing (solo)",
			args: args{
				artists: []spotify.SimpleArtist{{Name: "TestArtist1"}},
			},
			want: "TestArtist1",
		}, {
			name: "positive testing (empty)",
			args: args{
				artists: []spotify.SimpleArtist{},
			},
			want: "",
		}, {
			name: "positive testing (nil)",
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
