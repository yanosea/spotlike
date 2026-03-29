package spotlike

import (
	"context"
	"errors"
	"io"
	o "os"
	"strings"
	"testing"

	c "github.com/spf13/cobra"

	baseconfig "github.com/yanosea/spotlike/app/config"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/config"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/presenter"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewAuthCommand(t *testing.T) {
	output := ""
	exit := o.Exit
	origGetClientManagerFunc := api.GetClientManagerFunc

	type args struct {
		exit    func(int)
		cobra   proxy.Cobra
		version string
		conf    *config.SpotlikeCliConfig
		output  *string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				exit:    exit,
				cobra:   proxy.NewCobra(),
				version: "v0.0.0",
				conf: &config.SpotlikeCliConfig{
					SpotlikeConfig: baseconfig.SpotlikeConfig{
						SpotifyID:           "test_client_id",
						SpotifySecret:       "test_client_secret",
						SpotifyRedirectUri:  "test_redirect_uri",
						SpotifyRefreshToken: "test_refresh_token",
					},
				},
				output: &output,
			},
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			if tt.setup != nil {
				tt.setup(ctrl)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			got := NewAuthCommand(tt.args.exit, tt.args.cobra, tt.args.version, tt.args.conf, tt.args.output)
			if got == nil {
				t.Errorf("NewAuthCommand() = %v, want not nil", got)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run the auth command: %v", err)
				}
			}
		})
	}
}

func Test_runAuth(t *testing.T) {
	os := proxy.NewOs()
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()
	output := ""
	exit := o.Exit
	origGetClientManagerFunc := api.GetClientManagerFunc
	su := utility.NewStringsUtil()
	origPu := presenter.Pu
	origPrint := presenter.Print

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
			name: "negative testing (clientManager == nil)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: formatter.Red("❌ Client manager is not initialized..."),
			wantErr:    false,
			setup: func(_ *gomock.Controller) {
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
			name: "negative testing (clientManager.GetClient() failed)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err == nil {
						t.Errorf("Expected an error but got nil")
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				output = ""
			},
		},
		{
			name: "negative testing (client != nil)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: formatter.Yellow("⚡You've already setup your Spotify client..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				output = ""
			},
		},
		{
			name: "positive testing (ops.ID, ops.Secret, ops.RedirectUri are set)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!") + `
` + formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.") + `
  export SPOTIFY_ID=test_client_id
  export SPOTIFY_SECRET=test_client_secret
  export SPOTIFY_REDIRECT_URI=test_redirect_uri
  export SPOTIFY_REFRESH_TOKEN=test_refresh_token
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				output = ""
			},
		},
		{
			name: "negative testing (cancel input for client ID)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					exitCalled := false
					mockExit := func(code int) {
						exitCalled = true
						if code != 130 {
							t.Errorf("Exit code = %v, want %v", code, 130)
						}
					}
					authCmd := NewAuthCommand(
						mockExit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
					if !exitCalled {
						t.Error("Exit should have been called but wasn't")
					}
				},
			},
			wantStdOut: "\n" + formatter.Yellow("🚫 Cancelled authentication..."),
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				mockPromptui.EXPECT().NewPrompt().Return(mockPrompt)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\") failed in canceling input for client ID)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				mockPromptui.EXPECT().NewPrompt().Return(mockPrompt)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "\n") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, formatter.Yellow(\"🚫 Cancelled authentication...\")) failed in canceling input for client ID)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				mockPromptui.EXPECT().NewPrompt().Return(mockPrompt)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, formatter.Yellow("🚫 Cancelled authentication...")) {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.RunPrompt() failed in client ID input)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					err := authCmd.RunE(cmd, []string{})
					if err == nil || err.Error() != "prompt error" {
						t.Errorf("Expected 'prompt error' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt.EXPECT().Run().Return("", errors.New("prompt error"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				mockPromptui.EXPECT().NewPrompt().Return(mockPrompt)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "positive testing (with empty ID input)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!") + `
` + formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.") + `
  export SPOTIFY_ID=test_client_id
  export SPOTIFY_SECRET=test_client_secret
  export SPOTIFY_REDIRECT_URI=test_redirect_uri
  export SPOTIFY_REFRESH_TOKEN=test_refresh_token
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt2.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt3 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt3.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt3.EXPECT().SetMask('*')
				mockPrompt3.EXPECT().Run().Return("test_client_secret", nil)
				mockPrompt4 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt4.EXPECT().SetLabel("🔗 Input your Spotify Redirect URI")
				mockPrompt4.EXPECT().Run().Return("test_redirect_uri", nil)
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt3),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt4),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "negative testing (cancel input for client Secret)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					exitCalled := false
					exitCode := 0
					mockExit := func(code int) {
						exitCalled = true
						exitCode = code
					}
					authCmd := NewAuthCommand(
						mockExit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
					if !exitCalled {
						t.Error("Exit should have been called but wasn't")
					}
					if exitCode != 130 {
						t.Errorf("Exit code = %v, want %v", exitCode, 130)
					}
				},
			},
			wantStdOut: "\n" + formatter.Yellow("🚫 Cancelled authentication..."),
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\") in canceling input for client Secret)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "\n") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, formatter.Yellow(\"🚫 Cancelled authentication...\")) in canceling input for client Secret)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, formatter.Yellow("🚫 Cancelled authentication...")) {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.RunPromptWithMask() failed in client Secret input)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					err := authCmd.RunE(cmd, []string{})
					if err == nil || err.Error() != "failed to run prompt" {
						t.Errorf("Expected 'prompt error' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("", errors.New("failed to run prompt"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "positive testing (with empty Secret input",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!") + `
` + formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.") + `
  export SPOTIFY_ID=test_client_id
  export SPOTIFY_SECRET=test_client_secret
  export SPOTIFY_REDIRECT_URI=test_redirect_uri
  export SPOTIFY_REFRESH_TOKEN=test_refresh_token
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("", nil)
				mockPrompt3 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt3.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt3.EXPECT().SetMask('*')
				mockPrompt3.EXPECT().Run().Return("test_client_secret", nil)
				mockPrompt4 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt4.EXPECT().SetLabel("🔗 Input your Spotify Redirect URI")
				mockPrompt4.EXPECT().Run().Return("test_redirect_uri", nil)
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt3),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt4),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "negative testing (cancel input for Redirect URI)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					exitCalled := false
					exitCode := 0
					mockExit := func(code int) {
						exitCalled = true
						exitCode = code
					}
					authCmd := NewAuthCommand(
						mockExit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
					if !exitCalled {
						t.Error("Exit should have been called but wasn't")
					}
					if exitCode != 130 {
						t.Errorf("Exit code = %v, want %v", exitCode, 130)
					}
				},
			},
			wantStdOut: "\n" + formatter.Yellow("🚫 Cancelled authentication..."),
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("test_client_secret", nil)
				mockPrompt3 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt3.EXPECT().SetLabel("🔗 Input your Spotify Redirect URI")
				mockPrompt3.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt3),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\") in canceling input for Redirect URI)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("test_client_secret", nil)
				mockPrompt3 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt3.EXPECT().SetLabel("🔗 Input your Spotify Redirect URI")
				mockPrompt3.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt3),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "\n") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, formatter.Yellow(\"🚫 Cancelled authentication...\")) in canceling input for Redirect URI)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("test_client_secret", nil)
				mockPrompt3 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt3.EXPECT().SetLabel("🔗 Input your Spotify Redirect URI")
				mockPrompt3.EXPECT().Run().Return("", errors.New("^C"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt3),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, formatter.Yellow("🚫 Cancelled authentication...")) {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.RunPrompt() failed in Redirect URI input)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					err := authCmd.RunE(cmd, []string{})
					if err == nil || err.Error() != "failed to run prompt" {
						t.Errorf("Expected 'prompt error' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("test_client_secret", nil)
				mockPrompt3 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt3.EXPECT().SetLabel("🔗 Input your Spotify Redirect URI")
				mockPrompt3.EXPECT().Run().Return("", errors.New("failed to run prompt"))
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt3),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "positive testing (with empty Redirect URI input",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil {
						t.Errorf("Failed to run the auth command: %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!") + `
` + formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.") + `
  export SPOTIFY_ID=test_client_id
  export SPOTIFY_SECRET=test_client_secret
  export SPOTIFY_REDIRECT_URI=test_redirect_uri
  export SPOTIFY_REFRESH_TOKEN=test_refresh_token
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockPrompt1 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt1.EXPECT().SetLabel("🆔 Input your Spotify Client ID")
				mockPrompt1.EXPECT().Run().Return("test_client_id", nil)
				mockPrompt2 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt2.EXPECT().SetLabel("🔑 Input your Spotify Client Secret")
				mockPrompt2.EXPECT().SetMask('*')
				mockPrompt2.EXPECT().Run().Return("test_client_secret", nil)
				mockPrompt3 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt3.EXPECT().SetLabel("🔗 Input your Spotify Redirect URI")
				mockPrompt3.EXPECT().Run().Return("", nil)
				mockPrompt4 := proxy.NewMockPrompt(mockCtrl)
				mockPrompt4.EXPECT().SetLabel("🔗 Input your Spotify Redirect URI")
				mockPrompt4.EXPECT().Run().Return("test_redirect_uri", nil)
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				gomock.InOrder(
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt1),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt2),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt3),
					mockPromptui.EXPECT().NewPrompt().Return(mockPrompt4),
				)
				presenter.Pu = utility.NewPromptUtil(mockPromptui)
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Pu = origPu
				output = ""
			},
		},
		{
			name: "negative testing (clientManager.InitializeClient() failed)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:          "test_client_id",
								SpotifySecret:      "test_client_secret",
								SpotifyRedirectUri: "test_redirect_uri",
							},
						},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					err := authCmd.RunE(cmd, []string{})
					if err == nil {
						t.Errorf("Expected an error but got nil")
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
`,
			wantStdErr: "",
			wantOutput: formatter.Red("❌ Failed to authenticate... Please try again..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("failed to initialize client"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				output = ""
			},
		},
		{
			name: "negative testing (clientManager.GetClient() failed)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
						exit,
						proxy.NewCobra(),
						"v0.0.0",
						&config.SpotlikeCliConfig{
							SpotlikeConfig: baseconfig.SpotlikeConfig{
								SpotifyID:          "test_client_id",
								SpotifySecret:      "test_client_secret",
								SpotifyRedirectUri: "test_redirect_uri",
							},
						},
						&output,
					)
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					err := authCmd.RunE(cmd, []string{})
					if err == nil {
						t.Errorf("Expected an error but got nil")
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
`,
			wantStdErr: "",
			wantOutput: formatter.Red("❌ Failed to authenticate... Please try again..."),
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("failed to get client"))
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\") in displaying auth URL)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "\n") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"🌐 Login to Spotify by visiting the page below in your browser.\") in displaying auth URL)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "🌐 Login to Spotify by visiting the page below in your browser.") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"  \"+authUrl) in displaying auth URL)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "  "+"https://test-auth-url.com") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"\n\") in displaying auth URL (2nd time))",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				lfCount := 0
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "\n") {
						lfCount++
						if lfCount >= 2 {
							return errors.New("Print() failed")
						}
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (formatter.Green(\"🎉 Authentication succeeded!\"))",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, formatter.Green("🎉 Authentication succeeded!")) {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (formatter.Yellow(\"⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.\"))",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!"),
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.")) {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"  export SPOTIFY_ID=\"+conf.SpotifyID))",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!") + `
` + formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those."),
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "  export SPOTIFY_ID=test_client_id") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"  export SPOTIFY_SECRET=\"+conf.SpotifySecret))",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!") + `
` + formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.") + `
  export SPOTIFY_ID=test_client_id
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "  export SPOTIFY_SECRET=test_client_secret") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"  export SPOTIFY_REDIRECT_URI=\"+conf.SpotifyRedirectUri))",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!") + `
` + formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.") + `
  export SPOTIFY_ID=test_client_id
  export SPOTIFY_SECRET=test_client_secret
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "  export SPOTIFY_REDIRECT_URI=test_redirect_uri") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stdout, \"  export SPOTIFY_REFRESH_TOKEN=\"+conf.SpotifyRefreshToken))",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					authCmd := NewAuthCommand(
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
					cmd := &c.Command{}
					cmd.SetContext(context.Background())
					if err := authCmd.RunE(cmd, []string{}); err != nil && err.Error() != "Print() failed" {
						t.Errorf("Expected 'Print() failed' but got %v", err)
					}
				},
			},
			wantStdOut: `
🌐 Login to Spotify by visiting the page below in your browser.
  https://test-auth-url.com

` + formatter.Green("🎉 Authentication succeeded!") + `
` + formatter.Yellow("⚡ If you don't want spotlike to ask above again, execute commands below to set envs or set your profile to set those.") + `
  export SPOTIFY_ID=test_client_id
  export SPOTIFY_SECRET=test_client_secret
  export SPOTIFY_REDIRECT_URI=test_redirect_uri
`,
			wantStdErr: "",
			wantOutput: "",
			wantErr:    false,
			setup: func(mockCtrl *gomock.Controller) {
				mockClientManager := api.NewMockClientManager(mockCtrl)
				mockClient := api.NewMockClient(mockCtrl)
				mockClientManager.EXPECT().GetClient().Return(nil, errors.New("client not initialized"))
				mockClientManager.EXPECT().InitializeClient(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockClientManager.EXPECT().GetClient().Return(mockClient, nil)
				mockClient.EXPECT().Auth(
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(authUrlChan chan<- string, version string) (proxy.Client, string, error) {
					authUrlChan <- "https://test-auth-url.com"
					return nil, "test_refresh_token", nil
				})
				mockClient.EXPECT().UpdateConfig(gomock.Any()).Times(1)
				api.GetClientManagerFunc = func() api.ClientManager {
					return mockClientManager
				}
				presenter.Print = func(writer io.Writer, output string) error {
					if strings.Contains(output, "  export SPOTIFY_REFRESH_TOKEN=test_refresh_token") {
						return errors.New("Print() failed")
					}
					return origPrint(writer, output)
				}
			},
			cleanup: func() {
				api.GetClientManagerFunc = origGetClientManagerFunc
				presenter.Print = origPrint
				output = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			if tt.setup != nil {
				tt.setup(ctrl)
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
				t.Errorf("runAuth() gotStdOut doesn't match expected output")
			}
			cleanGotStdErr := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(gotStdErr)))
			cleanWantStdErr := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(tt.wantStdErr)))
			if cleanGotStdErr != cleanWantStdErr {
				t.Errorf("runAuth() gotStdErr = %v, want %v", cleanGotStdErr, cleanWantStdErr)
			}
			cleanOutput := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(output)))
			cleanWantOutput := su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(tt.wantOutput)))
			if cleanOutput != cleanWantOutput {
				t.Errorf("Output = %v, want %v", cleanOutput, cleanWantOutput)
			}
		})
	}
}
