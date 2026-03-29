package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/zmb3/spotify/v2"

	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"

	"github.com/yanosea/spotlike/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewArtistRepository(t *testing.T) {
	cm := api.NewClientManager(proxy.NewSpotify(), proxy.NewHttp(), proxy.NewRandstr(), proxy.NewUrl())

	tests := []struct {
		name string
		want artistDomain.ArtistRepository
	}{
		{
			name: "positive testing",
			want: &artistRepository{
				clientManager: cm,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArtistRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArtistRepository() = %v, want %v", got, tt.want)
			}
		})
	}
	if err := api.ResetClientManager(); err != nil {
		t.Errorf("Failed to reset client manager: %v", err)
	}
}

func Test_artistRepository_FindById(t *testing.T) {
	ma := &artistDomain.Artist{
		ID:   "test_artist_id",
		Name: "test_artist_name",
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
		want    *artistDomain.Artist
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
				mockApiClient.EXPECT().GetArtist(tt2.ctx, tt2.id).Return(&spotify.FullArtist{
					SimpleArtist: spotify.SimpleArtist{
						ID:   "test_artist_id",
						Name: "test_artist_name",
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
			name: "negative testing (client.GetArtist() failed)",
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
				mockApiClient.EXPECT().GetArtist(tt2.ctx, tt2.id).Return(nil, errors.New("failed to get artist"))
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
			r := &artistRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("artistRepository.FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("artistRepository.FindById() = nil, want %v", tt.want)
					return
				}
				if got.ID != tt.want.ID ||
					got.Name != tt.want.Name {
					t.Errorf("artistRepository.FindById() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_artistRepository_FindByNameLimit(t *testing.T) {
	ma := []*artistDomain.Artist{
		{
			ID:   "test_artist_id",
			Name: "test_artist_name",
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
		want    []*artistDomain.Artist
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
				name:  "test_artist_name",
				limit: 1,
			},
			want:    ma,
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields, tt2 *args) {
				mockApiClient := proxy.NewMockClient(mockCtrl)
				mockApiClient.EXPECT().Search(tt2.ctx, tt2.name, gomock.Any(), gomock.Any()).Return(&spotify.SearchResult{
					Artists: &spotify.FullArtistPage{
						Artists: []spotify.FullArtist{
							{
								SimpleArtist: spotify.SimpleArtist{
									ID:   "test_artist_id",
									Name: "test_artist_name",
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
				name:  "test_artist_name",
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
				name:  "test_artist_name",
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
			r := &artistRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.FindByNameLimit(tt.args.ctx, tt.args.name, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("artistRepository.FindByNameLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("artistRepository.FindByNameLimit() length = %v, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i].ID != tt.want[i].ID ||
						got[i].Name != tt.want[i].Name {
						t.Errorf("artistRepository.FindByNameLimit()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func Test_artistRepository_IsLiked(t *testing.T) {
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
				mockApiClient.EXPECT().CurrentUserFollows(tt2.ctx, "artist", tt2.id).Return([]bool{true}, nil)
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
			name: "negative testing (client.CurrentUserFollows() failed)",
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
				mockApiClient.EXPECT().CurrentUserFollows(tt2.ctx, "artist", tt2.id).Return(nil, errors.New("failed to get current user follows"))
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
			r := &artistRepository{
				clientManager: tt.fields.clientManager,
			}
			got, err := r.IsLiked(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("artistRepository.IsLiked() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("artistRepository.IsLiked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_artistRepository_Like(t *testing.T) {
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
				mockApiClient.EXPECT().FollowArtist(tt2.ctx, tt2.id).Return(nil)
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
			name: "negative testing (client.FollowArtist() failed)",
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
				mockApiClient.EXPECT().FollowArtist(tt2.ctx, tt2.id).Return(errors.New("failed to like artist"))
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
			r := &artistRepository{
				clientManager: tt.fields.clientManager,
			}
			err := r.Like(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("artistRepository.Like() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_artistRepository_Unlike(t *testing.T) {
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
				mockApiClient.EXPECT().UnfollowArtist(tt2.ctx, tt2.id).Return(nil)
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
			name: "negative testing (client.UnfollowArtist() failed)",
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
				mockApiClient.EXPECT().UnfollowArtist(tt2.ctx, tt2.id).Return(errors.New("failed to unlike artist"))
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
			r := &artistRepository{
				clientManager: tt.fields.clientManager,
			}
			err := r.Unlike(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("artistRepository.Unlike() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
