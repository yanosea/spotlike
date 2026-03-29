package spotlike

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"

	"go.uber.org/mock/gomock"
)

func TestNewGetAllAlbumsByArtistIdUseCase(t *testing.T) {
	type args struct {
		albumRepo albumDomain.AlbumRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *getAllAlbumsByArtistIdUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *getAllAlbumsByArtistIdUseCase
	}{
		{
			name: "positive testing",
			args: args{
				albumRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *getAllAlbumsByArtistIdUseCase {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				tt.albumRepo = mockAlbumRepo
				return &getAllAlbumsByArtistIdUseCase{
					albumRepo: mockAlbumRepo,
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
			if got := NewGetAllAlbumsByArtistIdUseCase(tt.args.albumRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetAllAlbumsByArtistIdUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAllAlbumsByArtistIdUseCase_Run(t *testing.T) {
	type fields struct {
		albumRepo albumDomain.AlbumRepository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*GetAllAlbumsByArtistIdUseCaseOutputDto
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				albumRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_artist_id",
			},
			want: []*GetAllAlbumsByArtistIdUseCaseOutputDto{
				{
					ID:          "test_album_id1",
					Name:        "test_album_name1",
					Artists:     "test_artist_name",
					ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          "test_album_id2",
					Name:        "test_album_name2",
					Artists:     "test_artist_name",
					ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				mockAlbumRepo.EXPECT().FindByArtistId(gomock.Any(), gomock.Any()).Return(
					[]*albumDomain.Album{
						{
							ID:   "test_album_id1",
							Name: "test_album_name1",
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_artist_id",
									Name: "test_artist_name",
								},
							},
							ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:   "test_album_id2",
							Name: "test_album_name2",
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_artist_id",
									Name: "test_artist_name",
								},
							},
							ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					nil,
				)
				tt.albumRepo = mockAlbumRepo
			},
		},
		{
			name: "negative testing (uc.albumRepo.FindByArtistId() failed)",
			fields: fields{
				albumRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_artist_id",
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				mockAlbumRepo.EXPECT().FindByArtistId(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get albums"))
				tt.albumRepo = mockAlbumRepo
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
			uc := &getAllAlbumsByArtistIdUseCase{
				albumRepo: tt.fields.albumRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllAlbumsByArtistIdUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAllAlbumsByArtistIdUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
