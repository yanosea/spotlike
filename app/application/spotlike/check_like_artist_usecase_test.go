package spotlike

import (
	"context"
	"reflect"
	"testing"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"

	"go.uber.org/mock/gomock"
)

func TestNewCheckLikeArtistUseCase(t *testing.T) {
	type args struct {
		artistRepo artistDomain.ArtistRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *checkLikeArtistUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *checkLikeArtistUseCase
	}{
		{
			name: "positive testing",
			args: args{
				artistRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *checkLikeArtistUseCase {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				tt.artistRepo = mockArtistRepo
				return &checkLikeArtistUseCase{
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
			if got := NewCheckLikeArtistUseCase(tt.args.artistRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCheckLikeArtistUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkLikeArtistUseCase_Run(t *testing.T) {
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
		want    bool
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
			want:    true,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				mockArtistRepo.EXPECT().IsLiked(gomock.Any(), gomock.Any()).Return(true, nil)
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
			uc := &checkLikeArtistUseCase{
				artistRepo: tt.fields.artistRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkLikeArtistUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkLikeArtistUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
