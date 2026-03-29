package spotlike

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/zmb3/spotify/v2"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"

	"go.uber.org/mock/gomock"
)

func TestNewGetAlbumUseCase(t *testing.T) {
	type args struct {
		albumRepo albumDomain.AlbumRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *getAlbumUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *getAlbumUseCase
	}{
		{
			name: "positive testing",
			args: args{
				albumRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *getAlbumUseCase {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				tt.albumRepo = mockAlbumRepo
				return &getAlbumUseCase{
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
			if got := NewGetAlbumUseCase(tt.args.albumRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetAlbumUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAlbumUseCase_Run(t *testing.T) {
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
		want    *GetAlbumUseCaseOutputDto
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
			want: &GetAlbumUseCaseOutputDto{
				ID:      "test_album_id",
				Name:    "test_album_name",
				Artists: "test_artist_name",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				mockAlbumRepo.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(
					&albumDomain.Album{
						ID:   "test_album_id",
						Name: "test_album_name",
						Artists: []spotify.SimpleArtist{
							{
								ID:   "test_artist_id",
								Name: "test_artist_name",
							},
						},
					},
					nil,
				)
				tt.albumRepo = mockAlbumRepo
			},
		},
		{
			name: "negative testing (uc.albumRepo.FindById() failed)",
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
				mockAlbumRepo.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get album"))
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
			uc := &getAlbumUseCase{
				albumRepo: tt.fields.albumRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAlbumUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAlbumUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
