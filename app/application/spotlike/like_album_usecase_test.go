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

func TestNewLikeAlbumUseCase(t *testing.T) {
	type args struct {
		albumDomain albumDomain.AlbumRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *likeAlbumUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *likeAlbumUseCase
	}{
		{
			name: "positive testing",
			args: args{
				albumDomain: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *likeAlbumUseCase {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				tt.albumDomain = mockAlbumRepo
				return &likeAlbumUseCase{
					albumDomain: mockAlbumRepo,
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
			if got := NewLikeAlbumUseCase(tt.args.albumDomain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLikeAlbumUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_likeAlbumUseCase_Run(t *testing.T) {
	type fields struct {
		albumDomain albumDomain.AlbumRepository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				albumDomain: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_album_id",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				mockAlbumRepo.EXPECT().Like(gomock.Any(), spotify.ID("test_album_id")).Return(nil)
				tt.albumDomain = mockAlbumRepo
			},
		},
		{
			name: "negative testing (uc.albumDomain.Like() failed)",
			fields: fields{
				albumDomain: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_album_id",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				mockAlbumRepo.EXPECT().Like(gomock.Any(), spotify.ID("test_album_id")).Return(errors.New("failed to like album"))
				tt.albumDomain = mockAlbumRepo
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
			uc := &likeAlbumUseCase{
				albumDomain: tt.fields.albumDomain,
			}
			if err := uc.Run(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("likeAlbumUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
