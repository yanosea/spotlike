package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"

	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"

	"github.com/yanosea/spotlike/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewTrackRepository(t *testing.T) {
	cm := api.NewClientManager(proxy.NewSpotify(), proxy.NewHttp(), proxy.NewRandstr(), proxy.NewUrl())

	tests := []struct {
		name string
		want trackDomain.TrackRepository
	}{
		{
			name: "positive testing",
			want: &trackRepository{
				clientManager: cm,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrackRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrackRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trackRepository_FindByArtistId(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02", "2000-01-01")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	mt := []*trackDomain.Track{
		{
			ID:   "test_track_id",
			Name: "test_track_name",
			Artists: []spotify.SimpleArtist{
				{
					ID:   "test_artist_id",
					Name: "test_artist_name",
				},
			},
			Album: spotify.SimpleAlbum{
				ID:   "test_album_id",
				Name: "test_album_name",
				Artists: []spotify.SimpleArtist{
					{
						ID:   "test_artist_id",
						Name: "test_artist_name",
					},
				},
			},
			TrackNumber: 1,
			ReleaseDate: expectedTime,
		},
	}

	type fields struct {
		clientManager api.ClientManager
	}
	type args struct {
		ctx context.Context
		id  spotify.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*trackDomain.Track
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields, tt2 *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test"),
			},
			want:    mt,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetArtistAlbums(tt2.ctx, tt2.id, nil).Return(&spotify.SimpleAlbumPage{
					Albums: []spotify.SimpleAlbum{
						{
							ID:   "test_album_id",
							Name: "test_album_name",
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_artist_id",
									Name: "test_artist_name",
								},
							},
							ReleaseDate:          "2000-01-01",
							ReleaseDatePrecision: "day",
						},
					},
				}, nil)
				mockApiClient.EXPECT().GetAlbumTracks(tt2.ctx, gomock.Any()).Return(&spotify.SimpleTrackPage{
					Tracks: []spotify.SimpleTrack{
						{
							ID:   "test_track_id",
							Name: "test_track_name",
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_artist_id",
									Name: "test_artist_name",
								},
							},
							TrackNumber: 1,
						},
					},
				}, nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (r.clientManager.GetClient() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test"),
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.GetArtistAlbums() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test"),
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetArtistAlbums(tt2.ctx, tt2.id, nil).Return(nil, errors.New("failed to get artist albums"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.GetAlbumTracks() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test"),
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetArtistAlbums(tt2.ctx, tt2.id, nil).Return(&spotify.SimpleAlbumPage{
					Albums: []spotify.SimpleAlbum{
						{
							ID:   "test_album_id",
							Name: "test_album_name",
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_album_artist_id",
									Name: "test_album_artist_name",
								},
							},
							ReleaseDate:          "2000-01-01",
							ReleaseDatePrecision: "day",
						},
					},
				}, nil)
				mockApiClient.EXPECT().GetAlbumTracks(tt2.ctx, gomock.Any()).Return(nil, errors.New("failed to get album tracks"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			r := &trackRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindByArtistId(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("trackRepository.FindByArtistId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("trackRepository.FindByArtistId() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID ||
					got[i].Name != tt.want[i].Name ||
					!reflect.DeepEqual(got[i].Artists, tt.want[i].Artists) ||
					!got[i].ReleaseDate.Equal(tt.want[i].ReleaseDate) {
					t.Errorf("trackRepository.FindByArtistId()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func Test_trackRepository_FindByAlbumId(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02", "2000-01-01")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	mt := []*trackDomain.Track{
		{
			ID:   "test_track_id",
			Name: "test_track_name",
			Artists: []spotify.SimpleArtist{
				{
					ID:   "test_artist_id",
					Name: "test_artist_name",
				},
			},
			Album: spotify.SimpleAlbum{
				ID:   "test_album_id",
				Name: "test_album_name",
				Artists: []spotify.SimpleArtist{
					{
						ID:   "test_artist_id",
						Name: "test_artist_name",
					},
				},
			},
			TrackNumber: 1,
			ReleaseDate: expectedTime,
		},
	}

	type fields struct {
		clientManager api.ClientManager
	}
	type args struct {
		ctx context.Context
		id  spotify.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*trackDomain.Track
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields, tt2 *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_album_id"),
			},
			want:    mt,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetAlbum(tt2.ctx, tt2.id).Return(&spotify.FullAlbum{
					SimpleAlbum: spotify.SimpleAlbum{
						ID:   "test_album_id",
						Name: "test_album_name",
						Artists: []spotify.SimpleArtist{
							{
								ID:   "test_artist_id",
								Name: "test_artist_name",
							},
						},
						ReleaseDate:          "2000-01-01",
						ReleaseDatePrecision: "day",
					},
				}, nil)
				mockApiClient.EXPECT().GetAlbumTracks(tt2.ctx, tt2.id).Return(&spotify.SimpleTrackPage{
					Tracks: []spotify.SimpleTrack{
						{
							ID:   "test_track_id",
							Name: "test_track_name",
							Artists: []spotify.SimpleArtist{
								{
									ID:   "test_artist_id",
									Name: "test_artist_name",
								},
							},
							TrackNumber: 1,
						},
					},
				}, nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (r.clientManager.GetClient() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_album_id"),
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.GetAlbum() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_album_id"),
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetAlbum(tt2.ctx, tt2.id).Return(nil, errors.New("failed to get album"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.GetAlbumTracks() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_album_id"),
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetAlbum(tt2.ctx, tt2.id).Return(&spotify.FullAlbum{
					SimpleAlbum: spotify.SimpleAlbum{
						ID:   "test_album_id",
						Name: "test_album_name",
						Artists: []spotify.SimpleArtist{
							{
								ID:   "test_artist_id",
								Name: "test_artist_name",
							},
						},
						ReleaseDate:          "2000-01-01",
						ReleaseDatePrecision: "day",
					},
				}, nil)
				mockApiClient.EXPECT().GetAlbumTracks(tt2.ctx, tt2.id).Return(nil, errors.New("failed to get album tracks"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			r := &trackRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindByAlbumId(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("trackRepository.FindByAlbumId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("trackRepository.FindByAlbumId() length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i].ID != tt.want[i].ID ||
						got[i].Name != tt.want[i].Name ||
						!reflect.DeepEqual(got[i].Artists, tt.want[i].Artists) ||
						!got[i].ReleaseDate.Equal(tt.want[i].ReleaseDate) {
						t.Errorf("trackRepository.FindByAlbumId()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func Test_trackRepository_FindById(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02", "2000-01-01")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	mt := &trackDomain.Track{
		ID:   "test_track_id",
		Name: "test_track_name",
		Artists: []spotify.SimpleArtist{
			{
				ID:   "test_artist_id",
				Name: "test_artist_name",
			},
		},
		Album: spotify.SimpleAlbum{
			ID:   "test_album_id",
			Name: "test_album_name",
			Artists: []spotify.SimpleArtist{
				{
					ID:   "test_artist_id",
					Name: "test_artist_name",
				},
			},
		},
		TrackNumber: 1,
		ReleaseDate: expectedTime,
	}

	type fields struct {
		clientManager api.ClientManager
	}
	type args struct {
		ctx context.Context
		id  spotify.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *trackDomain.Track
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields, tt2 *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			want:    mt,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetTrack(tt2.ctx, tt2.id).Return(&spotify.FullTrack{
					SimpleTrack: spotify.SimpleTrack{
						ID:   "test_track_id",
						Name: "test_track_name",
						Artists: []spotify.SimpleArtist{
							{
								ID:   "test_artist_id",
								Name: "test_artist_name",
							},
						},
						TrackNumber: 1,
					},
					Album: spotify.SimpleAlbum{
						ID:   "test_album_id",
						Name: "test_album_name",
						Artists: []spotify.SimpleArtist{
							{
								ID:   "test_artist_id",
								Name: "test_artist_name",
							},
						},
						ReleaseDate:          "2000-01-01",
						ReleaseDatePrecision: "day",
					},
				}, nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (r.clientManager.GetClient() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.GetTrack() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetTrack(tt2.ctx, tt2.id).Return(nil, errors.New("failed to get track"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			r := &trackRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("trackRepository.FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got == nil || tt.want == nil || got.ID != tt.want.ID ||
				got.Name != tt.want.Name ||
				!reflect.DeepEqual(got.Artists, tt.want.Artists) ||
				!got.ReleaseDate.Equal(tt.want.ReleaseDate)) {
				t.Errorf("trackRepository.FindById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trackRepository_FindByNameLimit(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02", "2000-01-01")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	mt := []*trackDomain.Track{
		{
			ID:   "test_track_id",
			Name: "test_track_name",
			Artists: []spotify.SimpleArtist{
				{
					ID:   "test_artist_id",
					Name: "test_artist_name",
				},
			},
			Album: spotify.SimpleAlbum{
				ID:   "test_album_id",
				Name: "test_album_name",
				Artists: []spotify.SimpleArtist{
					{
						ID:   "test_artist_id",
						Name: "test_artist_name",
					},
				},
			},
			TrackNumber: 1,
			ReleaseDate: expectedTime,
		},
	}

	type fields struct {
		clientManager api.ClientManager
	}
	type args struct {
		ctx   context.Context
		name  string
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*trackDomain.Track
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields, tt2 *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx:   context.Background(),
				name:  "test_track_name",
				limit: 1,
			},
			want:    mt,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().Search(tt2.ctx, tt2.name, gomock.Any(), gomock.Any()).Return(&spotify.SearchResult{
					Tracks: &spotify.FullTrackPage{
						Tracks: []spotify.FullTrack{
							{
								SimpleTrack: spotify.SimpleTrack{
									ID:   "test_track_id",
									Name: "test_track_name",
									Artists: []spotify.SimpleArtist{
										{
											ID:   "test_artist_id",
											Name: "test_artist_name",
										},
									},
									TrackNumber: 1,
								},
								Album: spotify.SimpleAlbum{
									ID:   "test_album_id",
									Name: "test_album_name",
									Artists: []spotify.SimpleArtist{
										{
											ID:   "test_artist_id",
											Name: "test_artist_name",
										},
									},
									ReleaseDate:          "2000-01-01",
									ReleaseDatePrecision: "day",
								},
							},
						},
					},
				}, nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (r.clientManager.GetClient() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx:   context.Background(),
				name:  "test_track_name",
				limit: 1,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.Search() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx:   context.Background(),
				name:  "test_track_name",
				limit: 1,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().Search(tt2.ctx, tt2.name, gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to search tracks"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			r := &trackRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindByNameLimit(tt.args.ctx, tt.args.name, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("trackRepository.FindByNameLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("trackRepository.FindByNameLimit() length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i].ID != tt.want[i].ID ||
						got[i].Name != tt.want[i].Name ||
						!reflect.DeepEqual(got[i].Artists, tt.want[i].Artists) ||
						!got[i].ReleaseDate.Equal(tt.want[i].ReleaseDate) {
						t.Errorf("trackRepository.FindByNameLimit()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func Test_trackRepository_IsLiked(t *testing.T) {
	type fields struct {
		clientManager api.ClientManager
	}
	type args struct {
		ctx context.Context
		id  spotify.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields, tt2 *args)
		cleanup func()
	}{
		{
			name: "positive testing (track is liked)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			want:    true,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().UserHasTracks(tt2.ctx, tt2.id).Return([]bool{true}, nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "positive testing (track is not liked)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			want:    false,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().UserHasTracks(tt2.ctx, tt2.id).Return([]bool{false}, nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (r.clientManager.GetClient() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			want:    false,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.UserHasTracks() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			want:    false,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().UserHasTracks(tt2.ctx, tt2.id).Return(nil, errors.New("failed to check if track is liked"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			r := &trackRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.IsLiked(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("trackRepository.IsLiked() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("trackRepository.IsLiked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trackRepository_Like(t *testing.T) {
	type fields struct {
		clientManager api.ClientManager
	}
	type args struct {
		ctx context.Context
		id  spotify.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields, tt2 *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().AddTracksToLibrary(tt2.ctx, tt2.id).Return(nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (r.clientManager.GetClient() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.AddTracksToLibrary() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().AddTracksToLibrary(tt2.ctx, tt2.id).Return(errors.New("failed to add track to library"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			r := &trackRepository{
				clientManager: tt.fields.clientManager,
			}
			if err := r.Like(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("trackRepository.Like() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_trackRepository_Unlike(t *testing.T) {
	type fields struct {
		clientManager api.ClientManager
	}
	type args struct {
		ctx context.Context
		id  spotify.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields, tt2 *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().RemoveTracksFromLibrary(tt2.ctx, tt2.id).Return(nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (r.clientManager.GetClient() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (client.RemoveTracksFromLibrary() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test_track_id"),
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().RemoveTracksFromLibrary(tt2.ctx, tt2.id).Return(errors.New("failed to remove track from library"))
				mockClient := api.NewMockClient(mockCtrl)
				mockClient.EXPECT().Open().Return(mockApiClient)
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				tt.clientManager = mockClientManager
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			r := &trackRepository{
				clientManager: tt.fields.clientManager,
			}
			if err := r.Unlike(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("trackRepository.Unlike() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
