package track

import (
	"reflect"
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"
)

func TestNewTrack(t *testing.T) {
	type args struct {
		id          spotify.ID
		name        string
		artists     []spotify.SimpleArtist
		album       spotify.SimpleAlbum
		trackNumber spotify.Numeric
		releaseDate time.Time
	}
	tests := []struct {
		name string
		args args
		want *Track
	}{
		{
			name: "positive testing",
			args: args{
				id:          "1",
				name:        "test track",
				artists:     []spotify.SimpleArtist{{ID: "1", Name: "test"}},
				album:       spotify.SimpleAlbum{ID: "1", Name: "test album"},
				trackNumber: 1,
				releaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: &Track{
				ID:          "1",
				Name:        "test track",
				Artists:     []spotify.SimpleArtist{{ID: "1", Name: "test"}},
				Album:       spotify.SimpleAlbum{ID: "1", Name: "test album"},
				TrackNumber: 1,
				ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrack(tt.args.id, tt.args.name, tt.args.artists, tt.args.album, tt.args.trackNumber, tt.args.releaseDate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrack() = %v, want %v", got, tt.want)
			}
		})
	}
}
