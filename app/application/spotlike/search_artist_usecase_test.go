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

func TestNewSearchArtistUseCase(t *testing.T) {
	type args struct {
		artistRepo artistDomain.ArtistRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *SearchArtistUseCaseStruct
		setup func(mockCtrl *gomock.Controller, tt *args) *SearchArtistUseCaseStruct
	}{
		{
			name: "positive testing",
			args: args{
				artistRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *SearchArtistUseCaseStruct {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				tt.artistRepo = mockArtistRepo
				return &SearchArtistUseCaseStruct{
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
			if got := NewSearchArtistUseCase(tt.args.artistRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSearchArtistUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchArtistUseCase_Run(t *testing.T) {
	type fields struct {
		artistRepo artistDomain.ArtistRepository
	}
	type args struct {
		ctx      context.Context
		keywords []string
		max      int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*SearchArtistUseCaseOutputDto
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields) []*SearchArtistUseCaseOutputDto
	}{
		{
			name: "positive testing",
			fields: fields{
				artistRepo: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "artist"},
				max:      5,
			},
			want:    nil,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) []*SearchArtistUseCaseOutputDto {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				artists := []*artistDomain.Artist{
					{
						ID:   spotify.ID("test_artist_id_1"),
						Name: "test_artist_name1",
					},
					{
						ID:   spotify.ID("test_artist_id_2"),
						Name: "test_artist_name2",
					},
				}
				mockArtistRepo.EXPECT().FindByNameLimit(
					gomock.Any(),
					"test artist",
					5,
				).Return(artists, nil)
				tt.artistRepo = mockArtistRepo
				expected := []*SearchArtistUseCaseOutputDto{
					{
						ID:   "test_artist_id_1",
						Name: "test_artist_name1",
					},
					{
						ID:   "test_artist_id_2",
						Name: "test_artist_name2",
					},
				}
				return expected
			},
		},
		{
			name: "negative testing (uc.artistRepo.FindByNameLimit() failed)",
			fields: fields{
				artistRepo: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "artist"},
				max:      5,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) []*SearchArtistUseCaseOutputDto {
				mockArtistRepo := artistDomain.NewMockArtistRepository(mockCtrl)
				var emptyArtists []*artistDomain.Artist
				mockArtistRepo.EXPECT().FindByNameLimit(
					gomock.Any(),
					"test artist",
					5,
				).Return(emptyArtists, errors.New("failed to find artists"))
				tt.artistRepo = mockArtistRepo
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.want = tt.setup(mockCtrl, &tt.fields)
			}
			uc := &SearchArtistUseCaseStruct{
				artistRepo: tt.fields.artistRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.keywords, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchArtistUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchArtistUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
