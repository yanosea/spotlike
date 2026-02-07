package formatter

import (
	"reflect"
	"testing"
	"time"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
)

func TestNewPlainFormatter(t *testing.T) {
	tests := []struct {
		name string
		want *PlainFormatter
	}{
		{
			name: "positive testing",
			want: &PlainFormatter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlainFormatter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlainFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlainFormatter_Format(t *testing.T) {
	type args struct {
		result any
	}
	tests := []struct {
		name    string
		f       *PlainFormatter
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive testing (result is GetVersionUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: &spotlikeApp.GetVersionUseCaseOutputDto{
					Version: "0.0.0",
				},
			},
			want:    "spotlike version 0.0.0",
			wantErr: false,
		},
		{
			name: "positive testing (result is SearchArtistUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: []*spotlikeApp.SearchArtistUseCaseOutputDto{
					{
						ID:   "artist_id_1",
						Name: "artist_name_1",
					},
					{
						ID:   "artist_id_2",
						Name: "artist_name_2",
					},
				},
			},
			want:    "[artist_id_1] Artist : artist_name_1\n[artist_id_2] Artist : artist_name_2",
			wantErr: false,
		},
		{
			name: "positive testing (result is SearchAlbumUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: []*spotlikeApp.SearchAlbumUseCaseOutputDto{
					{
						ID:          "album_id_1",
						Artists:     "artist_name_1",
						Name:        "album_name_1",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "album_id_2",
						Artists:     "artist_name_2",
						Name:        "album_name_2",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want:    "[album_id_1] Album : album_name_1 released at 2000-01-01 by artist_name_1\n[album_id_2] Album : album_name_2 released at 2000-01-01 by artist_name_2",
			wantErr: false,
		},
		{
			name: "positive testing (result is SearchTrackUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: []*spotlikeApp.SearchTrackUseCaseOutputDto{
					{
						ID:          "track_id_1",
						Artists:     "artist_name_1",
						Album:       "album_name_1",
						Name:        "track_name_1",
						TrackNumber: 1,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "track_id_2",
						Artists:     "artist_name_2",
						Album:       "album_name_2",
						Name:        "track_name_2",
						TrackNumber: 2,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want:    "[track_id_1] Track : #1 track_name_1 on album_name_1 released at 2000-01-01 by artist_name_1\n[track_id_2] Track : #2 track_name_2 on album_name_2 released at 2000-01-01 by artist_name_2",
			wantErr: false,
		},
		{
			name: "positive testing (result is GetArtistUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: []*spotlikeApp.GetArtistUseCaseOutputDto{
					{
						ID:   "artist_id_1",
						Name: "artist_name_1",
					},
					{
						ID:   "artist_id_2",
						Name: "artist_name_2",
					},
				},
			},
			want:    "[artist_id_1] Artist : artist_name_1\n[artist_id_2] Artist : artist_name_2",
			wantErr: false,
		},
		{
			name: "positive testing (result is GetAlbumUseCaseOutputDto array)",
			f:    &PlainFormatter{},
			args: args{
				result: []*spotlikeApp.GetAlbumUseCaseOutputDto{
					{
						ID:          "album_id_1",
						Name:        "album_name_1",
						Artists:     "artist_name_1",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "album_id_2",
						Name:        "album_name_2",
						Artists:     "artist_name_2",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want:    "[album_id_1] Album : album_name_1 released at 2000-01-01 by artist_name_1\n[album_id_2] Album : album_name_2 released at 2000-01-01 by artist_name_2",
			wantErr: false,
		},
		{
			name: "positive testing (result is GetTrackUseCaseOutputDto)",
			f:    &PlainFormatter{},
			args: args{
				result: []*spotlikeApp.GetTrackUseCaseOutputDto{
					{
						ID:          "track_id_1",
						Name:        "track_name_1",
						Artists:     "artist_name_1",
						Album:       "album_name_1",
						TrackNumber: 1,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "track_id_2",
						Name:        "track_name_2",
						Artists:     "artist_name_2",
						Album:       "album_name_2",
						TrackNumber: 2,
						ReleaseDate: time.Date(2001, 2, 2, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want:    "[track_id_1] Track : #1 track_name_1 on album_name_1 released at 2000-01-01 by artist_name_1\n[track_id_2] Track : #2 track_name_2 on album_name_2 released at 2001-02-02 by artist_name_2",
			wantErr: false,
		},
		{
			name: "negative testing (result is invalid)",
			f:    &PlainFormatter{},
			args: args{
				result: "invalid",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PlainFormatter{}
			got, err := f.Format(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlainFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PlainFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
