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

func TestNewGetAllTracksByAlbumIdUseCase(t *testing.T) {
	type args struct {
		trackRepo trackDomain.TrackRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *getAllTracksByAlbumIdUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *getAllTracksByAlbumIdUseCase
	}{
		{
			name: "positive testing",
			args: args{
				trackRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *getAllTracksByAlbumIdUseCase {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				tt.trackRepo = mockTrackRepo
				return &getAllTracksByAlbumIdUseCase{
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
			if got := NewGetAllTracksByAlbumIdUseCase(tt.args.trackRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetAllTracksByAlbumIdUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAllTracksByAlbumIdUseCase_Run(t *testing.T) {
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
		want    []*GetAllTracksByAlbumIdUseCaseOutputDto
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
				id:  "test_album_id",
			},
			want: []*GetAllTracksByAlbumIdUseCaseOutputDto{
				{
					ID:          "test_track_id1",
					Name:        "test_track_name1",
					Artists:     "test_artist_name",
					Album:       "test_album_name",
					TrackNumber: 1,
					ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          "test_track_id2",
					Name:        "test_track_name2",
					Artists:     "test_artist_name",
					Album:       "test_album_name",
					TrackNumber: 2,
					ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepo.EXPECT().FindByAlbumId(gomock.Any(), gomock.Any()).Return(
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
								ID:   "test_album_id",
								Name: "test_album_name",
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
								ID:   "test_album_id",
								Name: "test_album_name",
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
			name: "negative testing (uc.trackRepo.FindByAlbumId() failed)",
			fields: fields{
				trackRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_album_id",
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepo.EXPECT().FindByAlbumId(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get tracks"))
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
			uc := &getAllTracksByAlbumIdUseCase{
				trackRepo: tt.fields.trackRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllTracksByAlbumIdUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAllTracksByAlbumIdUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
