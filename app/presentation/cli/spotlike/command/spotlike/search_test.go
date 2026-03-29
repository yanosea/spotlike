package spotlike

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	c "github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
	baseconfig "github.com/yanosea/spotlike/app/config"
	albumDomain "github.com/yanosea/spotlike/app/domain/spotify/album"
	artistDomain "github.com/yanosea/spotlike/app/domain/spotify/artist"
	trackDomain "github.com/yanosea/spotlike/app/domain/spotify/track"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/config"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewSearchCommand(t *testing.T) {
	output := ""
	exit := os.Exit

	type args struct {
		cobra   proxy.Cobra
		authCmd proxy.Command
		output  *string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{
				cobra: proxy.NewCobra(),
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				output: &output,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSearchCommand(tt.args.cobra, tt.args.authCmd, tt.args.output)
			if got == nil {
				t.Errorf("NewSearchCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{"test", "artist"}); err != nil {
					t.Errorf("Failed to run the search command : %v", err)
				}
			}
		})
	}
}

func Test_runSearch(t *testing.T) {
	output := ""
	exit := os.Exit
	origSearchOps := searchOps
	origNewSearchArtistUseCase := spotlikeApp.NewSearchArtistUseCase
	origNewSearchAlbumUseCase := spotlikeApp.NewSearchAlbumUseCase
	origNewSearchTrackUseCase := spotlikeApp.NewSearchTrackUseCase
	origNewFormatter := formatter.NewFormatter
	su := utility.NewStringsUtil()

	type args struct {
		cmd     *c.Command
		authCmd proxy.Command
		output  *string
		args    []string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
		setup      func(mockCtrl *gomock.Controller, tt *args)
		cleanup    func()
	}{
		{
			name: "positive testing (search for an artist)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "artist"},
				output: &output,
			},
			wantOutput: "🆔ID🎤ARTISTtest_artist_idtest_artist_nameTOTAL:1artists!",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().Search(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&spotify.SearchResult{
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
					},
					nil,
				)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "positive testing (search for an album)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "album"},
				output: &output,
			},
			wantOutput: "🆔ID💿ALBUM🎤ARTISTS📅RELEASEDATEtest_album_idtest_album_nametest_artist_name0000-01-01TOTAL:1albums!",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Album = true
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().Search(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&spotify.SearchResult{
						Albums: &spotify.SimpleAlbumPage{
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
									ReleaseDate: "2000-01-01",
								},
							},
						},
					},
					nil,
				)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "positive testing (search for a track)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "track"},
				output: &output,
			},
			wantOutput: "🆔ID🔢NUMBER🎵TRACK💿ALBUM🎤ARTISTS📅RELEASEDATEtest_track_id1test_track_name0000-01-01TOTAL:1tracks!",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Track = true
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().Search(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&spotify.SearchResult{
						Tracks: &spotify.FullTrackPage{
							Tracks: []spotify.FullTrack{
								{
									SimpleTrack: spotify.SimpleTrack{
										ID:   "test_track_id",
										Name: "test_track_name",
										Album: spotify.SimpleAlbum{
											ID:          "test_album_id",
											Name:        "test_album_name",
											ReleaseDate: "2000-01-01",
											Artists: []spotify.SimpleArtist{
												{
													ID:   "test_artist_id",
													Name: "test_artist_name",
												},
											},
										},
										TrackNumber: 1,
									},
								},
							},
						},
					},
					nil,
				)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "negative testing (len(args) == 0)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{},
				output: &output,
			},
			wantOutput: formatter.Yellow("⚡Noqueryargumentsspecified..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "negative testing (!searchOps.Artist && !searchOps.Album && !searchOps.Track)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test"},
				output: &output,
			},
			wantOutput: formatter.Yellow("⚡Nosearchtypespecified..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = false
				searchOps.Album = false
				searchOps.Track = false
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "negative testing (searchOps.Artist && searchOps.Album || searchOps.Artist && searchOps.Track || searchOps.Album && searchOps.Track)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test"},
				output: &output,
			},
			wantOutput: formatter.Yellow("⚡Multiplesearchtypesisnotsupported..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				searchOps.Album = true
				searchOps.Track = true
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "negative testing (searchOps.Max < 1)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test"},
				output: &output,
			},
			wantOutput: formatter.Yellow("⚡Invalidsearchresultnumber..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				searchOps.Max = 0
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "negative testing (clientManager == nil)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test"},
				output: &output,
			},
			wantOutput: formatter.Red("❌Clientmanagerisnotinitialized..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				ctx := context.Background()
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "negative testing (err != nil && err.Error() == \"client not initialized\" and authCmd.RunE() failed)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test"},
				output: &output,
			},
			wantOutput: formatter.Red("❌Failedtoauthenticate...Pleasetryagain..."),
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				ctx := context.Background()
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().AuthURL(gomock.Any()).Return("https://example.com/auth")
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				mockUrl.EXPECT().Parse(gomock.Any()).Return(nil, errors.New("failed to parse"))
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockRandstr.EXPECT().Hex(11).Return("test-state")
				api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "negative testing (sAuc.Run((cmd.Context()), args, searchOps) failed)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "artist"},
				output: &output,
			},
			wantOutput: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				mockArtistRepository := artistDomain.NewMockArtistRepository(mockCtrl)
				mockArtistRepository.EXPECT().FindByNameLimit(ctx, "test artist", 10).Return(nil, errors.New("failed to find by name limit"))
				origNewSearchArtistUseCase := spotlikeApp.NewSearchArtistUseCase
				spotlikeApp.NewSearchArtistUseCase = func(artistRepo artistDomain.ArtistRepository) *spotlikeApp.SearchArtistUseCaseStruct {
					return origNewSearchArtistUseCase(mockArtistRepository)
				}
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				spotlikeApp.NewSearchArtistUseCase = origNewSearchArtistUseCase
				output = ""
			},
		},
		{
			name: "negative testing (len(sAoDtos) == 0)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "artist"},
				output: &output,
			},
			wantOutput: formatter.Yellow("⚡Noartistsfound..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				mockArtistRepository := artistDomain.NewMockArtistRepository(mockCtrl)
				mockArtistRepository.EXPECT().FindByNameLimit(ctx, "test artist", 10).Return([]*artistDomain.Artist{}, nil)
				spotlikeApp.NewSearchArtistUseCase = func(artistRepo artistDomain.ArtistRepository) *spotlikeApp.SearchArtistUseCaseStruct {
					return origNewSearchArtistUseCase(mockArtistRepository)
				}
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				spotlikeApp.NewSearchArtistUseCase = origNewSearchArtistUseCase
				output = ""
			},
		},
		{
			name: "negative testing (sauc.Run((cmd.Context()), args, searchOps) failed)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "album"},
				output: &output,
			},
			wantOutput: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Album = true
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				mockAlbumRepository := albumDomain.NewMockAlbumRepository(mockCtrl)
				mockAlbumRepository.EXPECT().FindByNameLimit(ctx, "test album", 10).Return(nil, errors.New("failed to find by name limit"))
				spotlikeApp.NewSearchAlbumUseCase = func(albumDomain.AlbumRepository) *spotlikeApp.SearchAlbumUseCaseStruct {
					return origNewSearchAlbumUseCase(mockAlbumRepository)
				}
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				spotlikeApp.NewSearchAlbumUseCase = origNewSearchAlbumUseCase
				output = ""
			},
		},
		{
			name: "negative testing (len(saoDtos) == 0)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "album"},
				output: &output,
			},
			wantOutput: formatter.Yellow("⚡Noalbumsfound..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Album = true
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				mockAlbumRepository := albumDomain.NewMockAlbumRepository(mockCtrl)
				mockAlbumRepository.EXPECT().FindByNameLimit(ctx, "test album", 10).Return([]*albumDomain.Album{}, nil)
				spotlikeApp.NewSearchAlbumUseCase = func(albumRepo albumDomain.AlbumRepository) *spotlikeApp.SearchAlbumUseCaseStruct {
					return origNewSearchAlbumUseCase(mockAlbumRepository)
				}
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				spotlikeApp.NewSearchAlbumUseCase = origNewSearchAlbumUseCase
				output = ""
			},
		},
		{
			name: "negative testing (stuc.Run((cmd.Context()), args, searchOps) failed)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "track"},
				output: &output,
			},
			wantOutput: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Track = true
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				mockTrackRepository := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepository.EXPECT().FindByNameLimit(ctx, "test track", 10).Return(nil, errors.New("failed to find by name limit"))
				spotlikeApp.NewSearchTrackUseCase = func(trackDomain.TrackRepository) *spotlikeApp.SearchTrackUseCaseStruct {
					return origNewSearchTrackUseCase(mockTrackRepository)
				}
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				spotlikeApp.NewSearchTrackUseCase = origNewSearchTrackUseCase
				output = ""
			},
		},
		{
			name: "negative testing (len(stoDtos) == 0)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "track"},
				output: &output,
			},
			wantOutput: formatter.Yellow("⚡Notracksfound..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Track = true
				ctx := context.Background()
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				mockTrackRepository := trackDomain.NewMockTrackRepository(mockCtrl)
				mockTrackRepository.EXPECT().FindByNameLimit(ctx, "test track", 10).Return([]*trackDomain.Track{}, nil)
				spotlikeApp.NewSearchTrackUseCase = func(trackDomain.TrackRepository) *spotlikeApp.SearchTrackUseCaseStruct {
					return origNewSearchTrackUseCase(mockTrackRepository)
				}
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				spotlikeApp.NewSearchTrackUseCase = origNewSearchTrackUseCase
				output = ""
			},
		},
		{
			name: "positive testing (formatter.NewFormatter(searchOps.Format) failed)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "artist"},
				output: &output,
			},
			wantOutput: formatter.Red("❌Failedtocreateaformatter..."),
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				searchOps.Format = "invalid"
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().Search(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&spotify.SearchResult{
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
					},
					nil,
				)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				searchOps = origSearchOps
				output = ""
			},
		},
		{
			name: "negative testing (f.Format() failed)",
			args: args{
				cmd: &c.Command{},
				authCmd: NewAuthCommand(
					exit,
					proxy.NewCobra(),
					"0.0.0",
					&config.SpotlikeCliConfig{
						SpotlikeConfig: baseconfig.SpotlikeConfig{
							SpotifyID:           "test_client_id",
							SpotifySecret:       "test_client_secret",
							SpotifyRedirectUri:  "test_redirect_uri",
							SpotifyRefreshToken: "test_refresh_token",
						},
					},
					&output,
				),
				args:   []string{"test", "artist"},
				output: &output,
			},
			wantOutput: "",
			wantErr:    true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				searchOps.Artist = true
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().Search(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&spotify.SearchResult{
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
					},
					nil,
				)
				mockAuthenticator := proxy.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().Client(gomock.Any(), gomock.Any()).Return(&http.Client{})
				mockSpotify := proxy.NewMockSpotify(mockCtrl)
				mockSpotify.EXPECT().NewAuthenticator(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAuthenticator)
				mockSpotify.EXPECT().NewClient(gomock.Any()).Return(mockSpotifyClient)
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockRandstr := proxy.NewMockRandstr(mockCtrl)
				mockUrl := proxy.NewMockUrl(mockCtrl)
				cm := api.NewClientManager(
					mockSpotify,
					mockHttp,
					mockRandstr,
					mockUrl,
				)
				if err := cm.InitializeClient(
					ctx,
					&api.ClientConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				mockFormatter := formatter.NewMockFormatter(mockCtrl)
				mockFormatter.EXPECT().Format(gomock.Any()).Return("", errors.New("format error"))
				formatter.NewFormatter = func(format string) (formatter.Formatter, error) {
					return mockFormatter, nil
				}
				cmd := &c.Command{}
				cmd.SetContext(ctx)
				tt.cmd = cmd
				output = ""
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				formatter.NewFormatter = origNewFormatter
				searchOps = origSearchOps
				output = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if err := runSearch(tt.args.cmd, tt.args.authCmd, tt.args.output, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
			*tt.args.output = su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(*tt.args.output)))
			if *tt.args.output != tt.wantOutput {
				t.Errorf("runSearch() got = %v, want %v", *tt.args.output, tt.wantOutput)
			}
		})
	}
}
