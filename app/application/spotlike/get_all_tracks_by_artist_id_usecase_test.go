package spotlike

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"

	"go.uber.org/mock/gomock"
)

func TestNewGetAllTracksByArtistIdUseCase(t *testing.T) {
	type args struct {
		tracksRepo trackDomain.TrackRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *getAllTracksByArtistIdUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *getAllTracksByArtistIdUseCase
	}{
		{
			name: "positive testing",
			args: args{
				tracksRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *getAllTracksByArtistIdUseCase {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				tt.tracksRepo = mockTrackRepo
				return &getAllTracksByArtistIdUseCase{
					trackRepo: mockTrackRepo,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.want = tt.setup(mockCtrl, &tt.args)
			}
			if got := NewGetAllTracksByArtistIdUseCase(tt.args.tracksRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetAllTracksByArtistIdUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAllTracksByArtistIdUseCase_Run(t *testing.T) {
	type fields struct {
		trackRepo trackDomain.TrackRepository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*GetAllTracksByArtistIdUseCaseOutputDto
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				trackRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_artist_id",
			},
			want: []*GetAllTracksByArtistIdUseCaseOutputDto{
				{
					ID:          "test_track_id1",
					Name:        "test_track_name1",
					Artists:     "test_artist_name",
					Album:       "test_album_name1",
					TrackNumber: 1,
					ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          "test_track_id2",
					Name:        "test_track_name2",
					Artists:     "test_artist_name",
					Album:       "test_album_name1",
					TrackNumber: 2,
					ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          "test_track_id3",
					Name:        "test_track_name3",
					Artists:     "test_artist_name",
					Album:       "test_album_name2",
					TrackNumber: 1,
					ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepo.EXPECT().FindByArtistId(gomock.Any(), gomock.Any()).Return(
					[]*trackDomain.Track{
						{
							ID:          "test_track_id1",
							Name:        "test_track_name1",
							TrackNumber: 1,
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_artist_id",
									Name: "test_artist_name",
								},
							},
							Album: spotify.SimpleAlbum{
								ID:   "test_album_id1",
								Name: "test_album_name1",
							},
							ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:          "test_track_id2",
							Name:        "test_track_name2",
							TrackNumber: 2,
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_artist_id",
									Name: "test_artist_name",
								},
							},
							Album: spotify.SimpleAlbum{
								ID:   "test_album_id1",
								Name: "test_album_name1",
							},
							ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:          "test_track_id3",
							Name:        "test_track_name3",
							TrackNumber: 1,
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_artist_id",
									Name: "test_artist_name",
								},
							},
							Album: spotify.SimpleAlbum{
								ID:   "test_album_id2",
								Name: "test_album_name2",
							},
							ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					nil,
				)
				tt.trackRepo = mockTrackRepo
			},
		},
		{
			name: "negative testing (uc.trackRepo.FindByArtistId() failed)",
			fields: fields{
				trackRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_artist_id",
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepo.EXPECT().FindByArtistId(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get tracks"))
				tt.trackRepo = mockTrackRepo
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			uc := &getAllTracksByArtistIdUseCase{
				trackRepo: tt.fields.trackRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllTracksByArtistIdUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("getAllTracksByArtistIdUseCase.Run() returned %d tracks, want %d tracks", len(got), len(tt.want))
					return
				}
				for i, track := range got {
					if track.ID != tt.want[i].ID {
						t.Errorf("Track[%d].ID = %v, want %v", i, track.ID, tt.want[i].ID)
					}
					if track.Name != tt.want[i].Name {
						t.Errorf("Track[%d].Name = %v, want %v", i, track.Name, tt.want[i].Name)
					}
					if track.Artists != tt.want[i].Artists {
						t.Errorf("Track[%d].Artists = %v, want %v", i, track.Artists, tt.want[i].Artists)
					}
					if track.Album != tt.want[i].Album {
						t.Errorf("Track[%d].Album = %v, want %v", i, track.Album, tt.want[i].Album)
					}
					if track.TrackNumber != tt.want[i].TrackNumber {
						t.Errorf("Track[%d].TrackNumber = %v, want %v", i, track.TrackNumber, tt.want[i].TrackNumber)
					}
					if !track.ReleaseDate.Equal(tt.want[i].ReleaseDate) {
						t.Errorf("Track[%d].ReleaseDate = %v, want %v", i, track.ReleaseDate, tt.want[i].ReleaseDate)
					}
				}
			}
		})
	}
}
