package spotlike

import (
	"context"
	"reflect"
	"testing"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"

	"go.uber.org/mock/gomock"
)

func TestNewCheckLikeAlbumUseCase(t *testing.T) {
	type args struct {
		albumRepo albumDomain.AlbumRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *checkLikeAlbumUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *checkLikeAlbumUseCase
	}{
		{
			name: "positive testing",
			args: args{
				albumRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *checkLikeAlbumUseCase {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				tt.albumRepo = mockAlbumRepo
				return &checkLikeAlbumUseCase{
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
			if got := NewCheckLikeAlbumUseCase(tt.args.albumRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCheckLikeAlbumUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkLikeAlbumUseCase_Run(t *testing.T) {
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
		want    bool
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
				id:  "test_album_id",
			},
			want:    true,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				mockAlbumRepo.EXPECT().IsLiked(gomock.Any(), gomock.Any()).Return(true, nil)
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
			uc := &checkLikeAlbumUseCase{
				albumRepo: tt.fields.albumRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkLikeAlbumUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkLikeAlbumUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
