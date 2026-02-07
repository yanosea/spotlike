package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"

	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"

	"github.com/yanosea/spotlike/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewAlbumRepository(t *testing.T) {
	cm := api.NewClientManager(proxy.NewSpotify(), proxy.NewHttp(), proxy.NewRandstr(), proxy.NewUrl())

	tests := []struct {
		name string
		want albumDomain.AlbumRepository
	}{
		{
			name: "positive testing",
			want: &albumRepository{
				clientManager: cm,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlbumRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlbumRepository() = %v, want %v", got, tt.want)
			}
		})
	}
	if err := api.ResetClientManager(); err != nil {
		t.Errorf("Failed to reset client manager: %v", err)
	}
}

func Test_albumRepository_FindByArtistId(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02", "2000-01-01")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	ma := []*albumDomain.Album{
		{
			ID:   "test_album_id",
			Name: "test_album_name",
			Artists: []spotify.SimpleArtist{
				{
					ID:   "test_album_artist_id",
					Name: "test_album_artist_name",
				},
			},
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
		want    []*albumDomain.Album
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
			want:    ma,
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
									ID:   "test_album_artist_id",
									Name: "test_album_artist_name",
								},
							},
							ReleaseDate:          "2000-01-01",
							ReleaseDatePrecision: "day",
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
			r := &albumRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindByArtistId(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("albumRepository.FindByArtistId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("albumRepository.FindByArtistId() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID ||
					got[i].Name != tt.want[i].Name ||
					!reflect.DeepEqual(got[i].Artists, tt.want[i].Artists) ||
					!got[i].ReleaseDate.Equal(tt.want[i].ReleaseDate) {
					t.Errorf("albumRepository.FindByArtistId()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func Test_albumRepository_FindById(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02", "2000-01-01")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	ma := &albumDomain.Album{
		ID:   "test_album_id",
		Name: "test_album_name",
		Artists: []spotify.SimpleArtist{
			{
				ID:   "test_album_artist_id",
				Name: "test_album_artist_name",
			},
		},
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
		want    *albumDomain.Album
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
			want:    ma,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().GetAlbum(tt2.ctx, tt2.id).Return(&spotify.FullAlbum{
					SimpleAlbum: spotify.SimpleAlbum{
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
			name: "negative testing (client.GetAlbum() failed)",
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
			r := &albumRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("albumRepository.FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("albumRepository.FindById() = nil, want %v", tt.want)
					return
				}
				if got.ID != tt.want.ID ||
					got.Name != tt.want.Name ||
					!reflect.DeepEqual(got.Artists, tt.want.Artists) ||
					!got.ReleaseDate.Equal(tt.want.ReleaseDate) {
					t.Errorf("albumRepository.FindById() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_albumRepository_FindByNameLimit(t *testing.T) {
	expectedTime, err := time.Parse("2006-01-02", "2000-01-01")
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	ma := []*albumDomain.Album{
		{
			ID:   "test_album_id",
			Name: "test_album_name",
			Artists: []spotify.SimpleArtist{
				{
					ID:   "test_album_artist_id",
					Name: "test_album_artist_name",
				},
			},
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
		want    []*albumDomain.Album
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
				name:  "test_album_name",
				limit: 1,
			},
			want:    ma,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().Search(tt2.ctx, tt2.name, gomock.Any(), gomock.Any()).Return(&spotify.SearchResult{
					Albums: &spotify.SimpleAlbumPage{
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
				name:  "test_album_name",
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
				name:  "test_album_name",
				limit: 1,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().Search(tt2.ctx, tt2.name, gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to search"))
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
			r := &albumRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindByNameLimit(tt.args.ctx, tt.args.name, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("albumRepository.FindByNameLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("albumRepository.FindByNameLimit() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID ||
					got[i].Name != tt.want[i].Name ||
					!reflect.DeepEqual(got[i].Artists, tt.want[i].Artists) ||
					!got[i].ReleaseDate.Equal(tt.want[i].ReleaseDate) {
					t.Errorf("albumRepository.FindByNameLimit()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func Test_albumRepository_IsLiked(t *testing.T) {
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
			name: "positive testing",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test"),
			},
			want:    true,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().UserHasAlbums(tt2.ctx, tt2.id).Return([]bool{true}, nil)
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
			name: "negative testing (client.UserHasAlbums() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test"),
			},
			want:    false,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().UserHasAlbums(tt2.ctx, tt2.id).Return(nil, errors.New("failed to get user has albums"))
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
			r := &albumRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.IsLiked(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("albumRepository.IsLiked() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("albumRepository.IsLiked() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_albumRepository_Like(t *testing.T) {
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
				id:  spotify.ID("test"),
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().AddAlbumsToLibrary(tt2.ctx, tt2.id).Return(nil)
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
			name: "negative testing (client.AddAlbumsToLibrary() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test"),
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().AddAlbumsToLibrary(tt2.ctx, tt2.id).Return(errors.New("failed to like album"))
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
			r := &albumRepository{
				clientManager: tt.fields.clientManager,
			}
			err := r.Like(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("albumRepository.Like() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_albumRepository_Unlike(t *testing.T) {
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
				id:  spotify.ID("test"),
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().RemoveAlbumsFromLibrary(tt2.ctx, tt2.id).Return(nil)
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
			name: "negative testing (client.RemoveAlbumsFromLibrary() failed)",
			fields: fields{
				clientManager: nil,
			},
			args: args{
				ctx: context.Background(),
				id:  spotify.ID("test"),
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().RemoveAlbumsFromLibrary(tt2.ctx, tt2.id).Return(errors.New("failed to unlike album"))
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
			r := &albumRepository{
				clientManager: tt.fields.clientManager,
			}
			err := r.Unlike(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("albumRepository.Unlike() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
