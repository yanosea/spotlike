package spotlike

import (
	"context"
	"errors"
	"reflect"
	"testing"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"

	"go.uber.org/mock/gomock"
)

func TestNewGetArtistUseCase(t *testing.T) {
	type args struct {
		artistRepo artistDomain.ArtistRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *GetArtistUseCaseStruct
		setup func(mockCtrl *gomock.Controller, tt *args) *GetArtistUseCaseStruct
	}{
		{
			name: "positive testing",
			args: args{
				artistRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *GetArtistUseCaseStruct {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				tt.artistRepo = mockArtistRepo
				return &GetArtistUseCaseStruct{
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
			if got := NewGetArtistUseCase(tt.args.artistRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetArtistUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getArtistUseCase_Run(t *testing.T) {
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
		want    *GetArtistUseCaseOutputDto
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
			want: &GetArtistUseCaseOutputDto{
				ID:   "test_artist_id",
				Name: "test_artist_name",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				mockArtistRepo.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(
					&artistDomain.Artist{
						ID:   "test_artist_id",
						Name: "test_artist_name",
					},
					nil,
				)
				tt.artistRepo = mockArtistRepo
			},
		},
		{
			name: "negative testing (uc.artistRepo.FindById() failed)",
			fields: fields{
				artistRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  "test_artist_id",
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				mockArtistRepo.EXPECT().FindById(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get artist"))
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
			uc := &GetArtistUseCaseStruct{
				artistRepo: tt.fields.artistRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getArtistUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getArtistUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
