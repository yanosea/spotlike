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

func TestNewUnlikeArtistUseCase(t *testing.T) {
	type args struct {
		artistRepo artistDomain.ArtistRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *unlikeArtistUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *unlikeArtistUseCase
	}{
		{
			name: "positive testing",
			args: args{
				artistRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *unlikeArtistUseCase {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				tt.artistRepo = mockArtistRepo
				return &unlikeArtistUseCase{
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
			if got := NewUnlikeArtistUseCase(tt.args.artistRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnlikeArtistUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unlikeArtistUseCase_Run(t *testing.T) {
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
				mockArtistRepo.EXPECT().Unlike(gomock.Any(), spotify.ID("test_artist_id")).Return(nil)
				tt.artistRepo = mockArtistRepo
			},
		},
		{
			name: "negative testing (uc.artistRepo.Unlike() failed)",
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
				mockArtistRepo.EXPECT().Unlike(gomock.Any(), spotify.ID("test_artist_id")).Return(errors.New("failed to unlike artist"))
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
			uc := &unlikeArtistUseCase{
				artistRepo: tt.fields.artistRepo,
			}
			if err := uc.Run(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("unlikeArtistUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
