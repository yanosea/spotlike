package spotlike

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/zmb3/spotify/v2"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"

	"go.uber.org/mock/gomock"
)

func TestNewLikeArtistUseCase(t *testing.T) {
	type args struct {
		artistRepo artistDomain.ArtistRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *likeArtistUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *likeArtistUseCase
	}{
		{
			name: "positive testing",
			args: args{
				artistRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *likeArtistUseCase {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				tt.artistRepo = mockArtistRepo
				return &likeArtistUseCase{
					artistRepo: mockArtistRepo,
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
			if got := NewLikeArtistUseCase(tt.args.artistRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLikeArtistUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_likeArtistUseCase_Run(t *testing.T) {
	type fields struct {
		artistRepo artistDomain.ArtistRepository
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
				artistRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_artist_id",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				mockArtistRepo.EXPECT().Like(gomock.Any(), spotify.ID("test_artist_id")).Return(nil)
				tt.artistRepo = mockArtistRepo
			},
		},
		{
			name: "negative testing (uc.artistRepo.Like() failed)",
			fields: fields{
				artistRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_artist_id",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				mockArtistRepo.EXPECT().Like(gomock.Any(), spotify.ID("test_artist_id")).Return(errors.New("failed to like artist"))
				tt.artistRepo = mockArtistRepo
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
			uc := &likeArtistUseCase{
				artistRepo: tt.fields.artistRepo,
			}
			if err := uc.Run(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("likeArtistUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
