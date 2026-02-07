package artist

import (
	"reflect"
	"testing"

	"github.com/zmb3/spotify/v2"
)

func TestNewArtist(t *testing.T) {
	type args struct {
		id   spotify.ID
		name string
	}
	tests := []struct {
		name string
		args args
		want *Artist
	}{
		{
			name: "positive testing",
			args: args{
				id:   "1",
				name: "artist",
			},
			want: &Artist{
				ID:   "1",
				Name: "artist",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArtist(tt.args.id, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArtist() = %v, want %v", got, tt.want)
			}
		})
	}
}
