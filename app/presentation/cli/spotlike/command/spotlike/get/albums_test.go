package get

import (
	"context"
	"errors"
	"net/http"
	o "os"
	"testing"

	c "github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"

	baseconfig "github.com/yanosea/spotlike/app/config"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command/spotlike"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/config"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/presenter"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewGetAlbumsCommand(t *testing.T) {
	output := ""
	exit := o.Exit

	type args struct {
		cobra   proxy.Cobra
		authCmd proxy.Command
		output  *string
	}
	tests := []struct {
		name string
		args args
		want proxy.Command
	}{
		{
			name: "positive testing",
			args: args{
				cobra: proxy.NewCobra(),
				authCmd: spotlike.NewAuthCommand(
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
			got := NewGetAlbumsCommand(tt.args.cobra, tt.args.authCmd, tt.args.output)
			if got == nil {
				t.Errorf("NewGetAlbumsCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{"test_artist_id"}); err != nil {
					t.Errorf("Failed to run the get albums command : %v", err)
				}
			}
		})
	}
}

func Test_runGetAlbums(t *testing.T) {
	os := proxy.NewOs()
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()
	output := ""
	exit := o.Exit
	origGetAlbumsOps := getAlbumsOps
	su := utility.NewStringsUtil()
	origPu := presenter.Pu
	origGetClientManagerFunc := api.GetClientManagerFunc
	origNewFormatter := formatter.NewFormatter

	type fields struct {
		Os        proxy.Os
		StdBuffer proxy.Buffer
		ErrBuffer proxy.Buffer
	}
	type args struct {
		fnc func()
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStdOut string
		wantStdErr string
		wantOutput string
		wantErr    bool
		setup      func(mockCtrl *gomock.Controller)
		cleanup    func()
	}{
		{
			name: "positive testing (normal case)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "test_client_id",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{"test_artist_id"}); err != nil {
						t.Errorf("Failed to run the getAlbums command: %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "\n🆔ID💿ALBUM🎤ARTISTS📅RELEASEDATEtest_album_idtest_album_nametest_artist_name2000-01-01TOTAL:1albums!",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().GetArtist(ctx, spotify.ID("test_artist_id")).Return(
					&spotify.FullArtist{
						SimpleArtist: spotify.SimpleArtist{
							ID:   "test_artist_id",
							Name: "test_artist_name",
						},
					},
					nil,
				)
				mockSpotifyClient.EXPECT().GetArtistAlbums(ctx, spotify.ID("test_artist_id"), nil).Return(
					&spotify.SimpleAlbumPage{
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
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				getAlbumsOps = origGetAlbumsOps
				output = ""
			},
		},
		{
			name: "negative testing (args is empty)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "test_client_id",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the getAlbums command: %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: formatter.Yellow("⚡ No ID arguments specified..."),
			wantErr:    false,
			setup:      nil,
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (clientManager is nil)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "test_client_id",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{"test_artist_id"}); err != nil {
						t.Errorf("Failed to run the getAlbums command: %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: formatter.Red("❌ Client manager is not initialized..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				api.GetClientManagerFunc = func() api.ClientManager {
					return nil
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				output = ""
			},
		},
		{
			name: "negative testing (client not initialized and authCmd.RunE() failed)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{"test_artist_id"}); err != nil {
						if err.Error() != "failed to run prompt" {
							t.Errorf("Failed to run the getAlbums command: %v", err)
						}
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized")).AnyTimes()
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt.EXPECT().Run().Return("", errors.New("failed to run prompt"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				mockPromptui.EXPECT().NewPrompt().Return(mockPrompt)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
			},
		},
		{
			name: "negative testing (failed to get artist)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "test_client_id",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{"test_artist_id"}); err != nil {
						if err.Error() != "failed to get artist" {
							t.Errorf("Failed to run the getAlbums command: %v", err)
						}
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().GetArtist(ctx, spotify.ID("test_artist_id")).Return(nil, errors.New("failed to get artist"))
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
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (gAucoDto is nil)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "test_client_id",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{"test_artist_id"}); err != nil {
						t.Errorf("Failed to run the getAlbums command: %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: formatter.Yellow("⚡ The id test_artist_id is not found... or it is not an artist..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().GetArtist(ctx, spotify.ID("test_artist_id")).Return(nil, errors.New("Resource not found"))
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
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				output = ""
			},
		},
		{
			name: "negative testing (failed to get all albums by artist id)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "test_client_id",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{"test_artist_id"}); err != nil {
						if err.Error() != "failed to get all albums by artist id" {
							t.Errorf("Failed to run the getAlbums command: %v", err)
						}
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().GetArtist(ctx, spotify.ID("test_artist_id")).Return(
					&spotify.FullArtist{
						SimpleArtist: spotify.SimpleArtist{
							ID:   "test_artist_id",
							Name: "test_artist_name",
						},
					},
					nil,
				)
				mockSpotifyClient.EXPECT().GetArtistAlbums(ctx, spotify.ID("test_artist_id"), nil).Return(nil, errors.New("failed to get all albums by artist id"))
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
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				output = ""
			},
		},
		{
			name: "negative testing (formatter.NewFormatter() failed)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "test_client_id",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					getAlbumsOps.Format = "test"
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{"test_artist_id"}); err != nil {
						if err.Error() != "invalid format" {
							t.Errorf("Failed to run the getAlbums command: %v", err)
						}
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: formatter.Red("❌ Failed to create a formatter..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().GetArtist(ctx, spotify.ID("test_artist_id")).Return(
					&spotify.FullArtist{
						SimpleArtist: spotify.SimpleArtist{
							ID:   "test_artist_id",
							Name: "test_artist_name",
						},
					},
					nil,
				)
				mockSpotifyClient.EXPECT().GetArtistAlbums(ctx, spotify.ID("test_artist_id"), nil).Return(
					&spotify.SimpleAlbumPage{
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
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				getAlbumsOps = origGetAlbumsOps
				output = ""
			},
		},
		{
			name: "negative testing (f.Format() failed)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := spotlike.NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:           "test_client_id",
								SpotifySecret:       "test_client_secret",
								SpotifyRedirectUri:  "test_redirect_uri",
								SpotifyRefreshToken: "test_refresh_token",
							},
						},
						&output,
					)
					getAlbumsCmd := NewGetAlbumsCommand(
						proxy.NewCobra(),
						authCmd,
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := getAlbumsCmd.RunE(cmd, []string{"test_artist_id"}); err != nil {
						if err.Error() != "format error" {
							t.Errorf("Failed to run the getAlbums command: %v", err)
						}
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				ctx := context.Background()
				mockSpotifyClient := proxy.NewMockClient(mockCtrl)
				mockSpotifyClient.EXPECT().GetArtist(ctx, spotify.ID("test_artist_id")).Return(
					&spotify.FullArtist{
						SimpleArtist: spotify.SimpleArtist{
							ID:   "test_artist_id",
							Name: "test_artist_name",
						},
					},
					nil,
				)
				mockSpotifyClient.EXPECT().GetArtistAlbums(ctx, spotify.ID("test_artist_id"), nil).Return(
					&spotify.SimpleAlbumPage{
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
			},
			cleanup: func() {
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				formatter.NewFormatter = origNewFormatter
				getAlbumsOps = origGetAlbumsOps
				output = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			c := utility.NewCapturer(tt.fields.Os, tt.fields.StdBuffer, tt.fields.ErrBuffer)
			gotStdOut, gotStdErr, err := c.CaptureOutput(tt.args.fnc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			cleanGotStdOut := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(gotStdOut)))
			cleanWantStdOut := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(tt.wantStdOut)))
			if cleanGotStdOut != cleanWantStdOut {
				t.Logf("gotStdOut: %v", gotStdOut)
				t.Logf("wantStdOut: %v", tt.wantStdOut)
				t.Errorf("runGetAlbums() gotStdOut doesn't match expected output")
			}
			cleanGotStdErr := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(gotStdErr)))
			cleanWantStdErr := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(tt.wantStdErr)))
			if cleanGotStdErr != cleanWantStdErr {
				t.Errorf("runGetAlbums() gotStdErr = %v, want %v", cleanGotStdErr, cleanWantStdErr)
			}
			cleanOutput := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(output)))
			cleanWantOutput := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(tt.wantOutput)))
			if cleanOutput != cleanWantOutput {
				t.Errorf("Output = %v, want %v", cleanOutput, cleanWantOutput)
			}
		})
	}
}
