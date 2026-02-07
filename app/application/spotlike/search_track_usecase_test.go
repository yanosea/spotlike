package spotlike

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"

	"go.uber.org/mock/gomock"
)

func TestNewSearchTrackUseCase(t *testing.T) {
	type args struct {
		trackRepo trackDomain.TrackRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *SearchTrackUseCaseStruct
		setup func(mockCtrl *gomock.Controller, tt *args) *SearchTrackUseCaseStruct
	}{
		{
			name: "positive testing",
			args: args{
				trackRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *SearchTrackUseCaseStruct {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				tt.trackRepo = mockTrackRepo
				return &SearchTrackUseCaseStruct{
					trackRepo: mockTrackRepo,
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
			if got := NewSearchTrackUseCase(tt.args.trackRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSearchTrackUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchTrackUseCase_Run(t *testing.T) {
	type fields struct {
		trackRepo trackDomain.TrackRepository
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
		want    []*SearchTrackUseCaseOutputDto
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields) []*SearchTrackUseCaseOutputDto
	}{
		{
			name: "positive testing",
			fields: fields{
				trackRepo: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "track"},
				max:      5,
			},
			want:    nil,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) []*SearchTrackUseCaseOutputDto {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				tracks := []*trackDomain.Track{
					{
						ID:   spotify.ID("test_track_id_1"),
						Name: "test_track_name1",
						Artists: []spotify.SimpleArtist{
							{Name: "test_artist_name1"},
							{Name: "test_artist_name2"},
						},
						Album: spotify.SimpleAlbum{
							Name:                 "test_album_name1",
							ReleaseDate:          "2000-01-01",
							ReleaseDatePrecision: "day",
						},
						TrackNumber: 1,
					},
					{
						ID:   spotify.ID("test_track_id_2"),
						Name: "test_track_name2",
						Artists: []spotify.SimpleArtist{
							{Name: "test_artist_name3"},
						},
						Album: spotify.SimpleAlbum{
							Name:                 "test_album_name2",
							ReleaseDate:          "2000-01-01",
							ReleaseDatePrecision: "day",
						},
						TrackNumber: 2,
					},
				}
				mockTrackRepo.EXPECT().FindByNameLimit(
					gomock.Any(),
					"test track",
					5,
				).Return(tracks, nil)
				tt.trackRepo = mockTrackRepo
				expected := []*SearchTrackUseCaseOutputDto{
					{
						ID:          "test_track_id_1",
						Artists:     "test_artist_name1, test_artist_name2",
						Album:       "test_album_name1",
						Name:        "test_track_name1",
						TrackNumber: 1,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          "test_track_id_2",
						Artists:     "test_artist_name3",
						Album:       "test_album_name2",
						Name:        "test_track_name2",
						TrackNumber: 2,
						ReleaseDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				return expected
			},
		},
		{
			name: "negative testing (uc.trackRepo.FindByNameLimit() failed)",
			fields: fields{
				trackRepo: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "track"},
				max:      5,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) []*SearchTrackUseCaseOutputDto {
				mockTrackRepo := trackDomain.NewMockTrackRepository(mockCtrl)
				var emptyTracks []*trackDomain.Track
				mockTrackRepo.EXPECT().FindByNameLimit(
					gomock.Any(),
					"test track",
					5,
				).Return(emptyTracks, errors.New("failed to find tracks"))
				tt.trackRepo = mockTrackRepo
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
			uc := &SearchTrackUseCaseStruct{
				trackRepo: tt.fields.trackRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.keywords, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchTrackUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchTrackUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
