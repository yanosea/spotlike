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

func TestNewSearchAlbumUseCase(t *testing.T) {
	type args struct {
		albumRepo albumDomain.AlbumRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *SearchAlbumUseCaseStruct
		setup func(mockCtrl *gomock.Controller, tt *args) *SearchAlbumUseCaseStruct
	}{
		{
			name: "positive testing",
			args: args{
				albumRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *SearchAlbumUseCaseStruct {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				tt.albumRepo = mockAlbumRepo
				return &SearchAlbumUseCaseStruct{
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
			if got := NewSearchAlbumUseCase(tt.args.albumRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSearchAlbumUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchAlbumUseCase_Run(t *testing.T) {
	type fields struct {
		albumRepo albumDomain.AlbumRepository
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
		want    []*SearchAlbumUseCaseOutputDto
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields) []*SearchAlbumUseCaseOutputDto
	}{
		{
			name: "positive testing",
			fields: fields{
				albumRepo: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "album"},
				max:      5,
			},
			want:    nil,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) []*SearchAlbumUseCaseOutputDto {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				albums := []*albumDomain.Album{
					{
						ID:   spotify.ID("test_album_id_1"),
						Name: "test_album_name",
						Artists: []spotify.SimpleArtist{
							{Name: "test_artist_name1"},
							{Name: "test_artist_name2"},
						},
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:   spotify.ID("test_album_id_2"),
						Name: "test_album_name3",
						Artists: []spotify.SimpleArtist{
							{Name: "test_artist_name3"},
						},
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				mockAlbumRepo.EXPECT().FindByNameLimit(
					gomock.Any(),
					"test album",
					5,
				).Return(albums, nil)
				tt.albumRepo = mockAlbumRepo
				expected := []*SearchAlbumUseCaseOutputDto{
					{
						ID:          "test_album_id_1",
						Artists:     "test_artist_name1, test_artist_name2",
						Name:        "test_album_name",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "test_album_id_2",
						Artists:     "test_artist_name3",
						Name:        "test_album_name3",
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				return expected
			},
		},
		{
			name: "negative testing (uc.albumRepo.FindByNameLimit() failed)",
			fields: fields{
				albumRepo: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "album"},
				max:      5,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) []*SearchAlbumUseCaseOutputDto {
				mockAlbumRepo := albumDomain.NewMockAlbumRepository(mockCtrl)
				var emptyAlbums []*albumDomain.Album
				mockAlbumRepo.EXPECT().FindByNameLimit(
					gomock.Any(),
					"test album",
					5,
				).Return(emptyAlbums, errors.New("failed to find albums"))
				tt.albumRepo = mockAlbumRepo
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
			uc := &SearchAlbumUseCaseStruct{
				albumRepo: tt.fields.albumRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.keywords, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchAlbumUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchAlbumUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
