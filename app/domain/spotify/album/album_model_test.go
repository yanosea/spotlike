package album

import (
	"reflect"
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"
)

func TestNewAlbum(t *testing.T) {
	type args struct {
		id          spotify.ID
		name        string
		artists     []spotify.SimpleArtist
		releaseDate time.Time
	}
	tests := []struct {
		name string
		args args
		want *Album
	}{
		{
			name: "positive testing",
			args: args{
				id:          "1",
				name:        "album",
				artists:     []spotify.SimpleArtist{{ID: "1", Name: "artist"}},
				releaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: &Album{
				ID:          "1",
				Name:        "album",
				Artists:     []spotify.SimpleArtist{{ID: "1", Name: "artist"}},
				ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlbum(tt.args.id, tt.args.name, tt.args.artists, tt.args.releaseDate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlbum() = %v, want %v", got, tt.want)
			}
		})
	}
}
