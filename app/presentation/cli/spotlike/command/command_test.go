package command

import (
	"context"
	"errors"
	"io"
	o "os"
	"strings"
	"testing"

	"github.com/fatih/color"

	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/presenter"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"
	"go.uber.org/mock/gomock"
)

func Test_newCli(t *testing.T) {
	cobra := proxy.NewCobra()
	ctx := context.Background()
	exit := func(code int) {}

	type args struct {
		exit  func(int)
		cobra proxy.Cobra
		ctx   context.Context
	}
	tests := []struct {
		name string
		args args
		want *cli
	}{
		{
			name: "positive testing",
			args: args{
				exit:  exit,
				cobra: cobra,
				ctx:   ctx,
			},
			want: &cli{
				Exit:          exit,
				Cobra:         cobra,
				RootCommand:   nil,
				Context:       ctx,
				ClientManager: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCli(tt.args.exit, tt.args.cobra, tt.args.ctx).(*cli)
			if got.Cobra != tt.want.Cobra {
				t.Errorf("newCli().Cobra = %v, want %v", got.Cobra, tt.want.Cobra)
			}
			if got.RootCommand != tt.want.RootCommand {
				t.Errorf("newCli().RootCommand = %v, want %v", got.RootCommand, tt.want.RootCommand)
			}
			if got.Context != tt.want.Context {
				t.Errorf("newCli().Context = %v, want %v", got.Context, tt.want.Context)
			}
			if got.ClientManager != tt.want.ClientManager {
				t.Errorf("newCli().ClientManager = %v, want %v", got.ClientManager, tt.want.ClientManager)
			}
		})
	}
}

func Test_cli_Init(t *testing.T) {
	osProxy := proxy.NewOs()
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()
	cobra := proxy.NewCobra()
	ctx := context.Background()
	exit := func(code int) {}
	origPrint := presenter.Print

	type fields struct {
		os        proxy.Os
		StdBuffer proxy.Buffer
		ErrBuffer proxy.Buffer
	}
	type args struct {
		fnc func(mockCtrl *gomock.Controller)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStdOut string
		wantStdErr string
		wantErr    bool
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(_ *gomock.Controller) {
					c := &cli{
						Exit:          exit,
						Cobra:         cobra,
						RootCommand:   nil,
						Context:       ctx,
						ClientManager: nil,
					}
					if got := c.Init(
						proxy.NewEnvconfig(),
						proxy.NewSpotify(),
						proxy.NewHttp(),
						proxy.NewRandstr(),
						proxy.NewUrl(),
						"0.0.0",
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 0 {
						t.Errorf("cli.Init() = %v, want %v", got, 0)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (configurator.GetConfig() failed)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
					mockEnvconfig.EXPECT().Process("", gomock.Any()).Return(errors.New("EnvconfigProxy.Process() failed"))
					c := &cli{
						Exit:          exit,
						Cobra:         cobra,
						RootCommand:   nil,
						Context:       ctx,
						ClientManager: nil,
					}
					if got := c.Init(
						mockEnvconfig,
						proxy.NewSpotify(),
						proxy.NewHttp(),
						proxy.NewRandstr(),
						proxy.NewUrl(),
						"0.0.0",
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 1 {
						t.Errorf("cli.Init() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: color.RedString("Error : EnvconfigProxy.Process() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stderr, output) failed in Init)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					defer func() { presenter.Print = origPrint }()
					mockEnvconfig := proxy.NewMockEnvconfig(mockCtrl)
					mockEnvconfig.EXPECT().Process("", gomock.Any()).Return(errors.New("EnvconfigProxy.Process() failed"))
					presenter.Print = func(writer io.Writer, output string) error {
						if writer == o.Stderr && strings.Contains(output, "EnvconfigProxy.Process() failed") {
							return errors.New("Print() failed")
						}
						return origPrint(writer, output)
					}

					c := &cli{
						Exit:          exit,
						Cobra:         cobra,
						RootCommand:   nil,
						Context:       ctx,
						ClientManager: nil,
					}
					if got := c.Init(
						mockEnvconfig,
						proxy.NewSpotify(),
						proxy.NewHttp(),
						proxy.NewRandstr(),
						proxy.NewUrl(),
						"0.0.0",
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 1 {
						t.Errorf("cli.Init() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (ClientManager.InitializeClient() failed)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().InitializeClient(gomock.Any(), gomock.Any()).Return(errors.New("ClientManager.InitializeClient() failed"))
					c := &cli{
						Exit:          exit,
						Cobra:         cobra,
						RootCommand:   nil,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Init(
						proxy.NewEnvconfig(),
						proxy.NewSpotify(),
						proxy.NewHttp(),
						proxy.NewRandstr(),
						proxy.NewUrl(),
						"0.0.0",
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 1 {
						t.Errorf("cli.Init() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: color.RedString("Error : ClientManager.InitializeClient() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
				if err := o.Setenv("SPOTIFY_ID", "test_id"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("SPOTIFY_SECRET", "test_secret"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/callback"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("SPOTIFY_REFRESH_TOKEN", "test_refresh_token"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
			},
			cleanup: func() {
				output = ""
				if err := o.Unsetenv("SPOTIFY_ID"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("SPOTIFY_SECRET"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("SPOTIFY_REDIRECT_URI"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("SPOTIFY_REFRESH_TOKEN"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stderr, output) failed in InitializeClient)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					defer func() { presenter.Print = origPrint }()
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().InitializeClient(gomock.Any(), gomock.Any()).Return(errors.New("ClientManager.InitializeClient() failed"))
					presenter.Print = func(writer io.Writer, output string) error {
						if writer == o.Stderr && strings.Contains(output, "ClientManager.InitializeClient() failed") {
							return errors.New("Print() failed")
						}
						return origPrint(writer, output)
					}
					c := &cli{
						Exit:          exit,
						Cobra:         cobra,
						RootCommand:   nil,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Init(
						proxy.NewEnvconfig(),
						proxy.NewSpotify(),
						proxy.NewHttp(),
						proxy.NewRandstr(),
						proxy.NewUrl(),
						"0.0.0",
						utility.NewVersionUtil(
							proxy.NewDebug(),
						),
					); got != 1 {
						t.Errorf("cli.Init() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
				if err := o.Setenv("SPOTIFY_ID", "test_id"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("SPOTIFY_SECRET", "test_secret"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/callback"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
				if err := o.Setenv("SPOTIFY_REFRESH_TOKEN", "test_refresh_token"); err != nil {
					t.Errorf("Failed to set environment variable: %v", err)
				}
			},
			cleanup: func() {
				output = ""
				if err := o.Unsetenv("SPOTIFY_ID"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("SPOTIFY_SECRET"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("SPOTIFY_REDIRECT_URI"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
				if err := o.Unsetenv("SPOTIFY_REFRESH_TOKEN"); err != nil {
					t.Errorf("Failed to unset environment variable: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			output = ""
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			c := utility.NewCapturer(tt.fields.os, tt.fields.StdBuffer, tt.fields.ErrBuffer)
			gotStdOut, gotStdErr, err := c.CaptureOutput(func() {
				tt.args.fnc(mockCtrl)
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStdOut != tt.wantStdOut {
				t.Errorf("Init() gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("Init() gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}

func Test_cli_Run(t *testing.T) {
	osProxy := proxy.NewOs()
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()
	ctx := context.Background()
	exit := func(code int) {}
	origPrint := presenter.Print

	type fields struct {
		os        proxy.Os
		StdBuffer proxy.Buffer
		ErrBuffer proxy.Buffer
	}
	type args struct {
		fnc func(mockCtrl *gomock.Controller)
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStdOut string
		wantStdErr string
		wantErr    bool
		setup      func()
		cleanup    func()
	}{
		{
			name: "positive testing",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().IsClientInitialized().Return(true)
					mockClientManager.EXPECT().CloseClient().Return(nil)
					c := &cli{
						Exit:          exit,
						Cobra:         proxy.NewCobra(),
						RootCommand:   mockCommand,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Run(); got != 0 {
						t.Errorf("cli.Run() = %v, want %v", got, 0)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "positive testing (no client initialized)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().IsClientInitialized().Return(false)
					c := &cli{
						Exit:          exit,
						Cobra:         proxy.NewCobra(),
						RootCommand:   mockCommand,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Run(); got != 0 {
						t.Errorf("cli.Run() = %v, want %v", got, 0)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "positive testing (no client manager)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					c := &cli{
						Exit:          exit,
						Cobra:         proxy.NewCobra(),
						RootCommand:   mockCommand,
						Context:       ctx,
						ClientManager: nil,
					}
					if got := c.Run(); got != 0 {
						t.Errorf("cli.Run() = %v, want %v", got, 0)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "positive testing (with output)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().IsClientInitialized().Return(true)
					mockClientManager.EXPECT().CloseClient().Return(nil)
					c := &cli{
						Exit:          exit,
						Cobra:         proxy.NewCobra(),
						RootCommand:   mockCommand,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Run(); got != 0 {
						t.Errorf("cli.Run() = %v, want %v", got, 0)
					}
				},
			},
			wantStdOut: "test output\n",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = "test output"
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (c.RootCommand.ExecuteContext() failed)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(errors.New("CommandProxy.ExecuteContext() failed"))
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().IsClientInitialized().Return(true)
					mockClientManager.EXPECT().CloseClient().Return(nil)
					c := &cli{
						Exit:          exit,
						Cobra:         proxy.NewCobra(),
						RootCommand:   mockCommand,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Run(); got != 1 {
						t.Errorf("cli.Run() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: color.RedString("Error : CommandProxy.ExecuteContext() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(out, output) failed in Run)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					defer func() { presenter.Print = origPrint }()
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().IsClientInitialized().Return(true)
					mockClientManager.EXPECT().CloseClient().Return(nil)
					presenter.Print = func(writer io.Writer, output string) error {
						return errors.New("Print() failed")
					}
					c := &cli{
						Exit:          exit,
						Cobra:         proxy.NewCobra(),
						RootCommand:   mockCommand,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Run(); got != 1 {
						t.Errorf("cli.Run() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = "test output"
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (c.ClientManager.CloseClient() failed)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().IsClientInitialized().Return(true)
					mockClientManager.EXPECT().CloseClient().Return(errors.New("ClientManager.CloseClient() failed"))
					c := &cli{
						Exit:          exit,
						Cobra:         proxy.NewCobra(),
						RootCommand:   mockCommand,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Run(); got != 1 {
						t.Errorf("cli.Run() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: color.RedString("Error : ClientManager.CloseClient() failed") + "\n",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
		{
			name: "negative testing (presenter.Print(os.Stderr, output) failed in CloseClient error)",
			fields: fields{
				os:        osProxy,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func(mockCtrl *gomock.Controller) {
					defer func() { presenter.Print = origPrint }()
					mockCommand := proxy.NewMockCommand(mockCtrl)
					mockCommand.EXPECT().ExecuteContext(gomock.Any()).Return(nil)
					mockClientManager := api.NewMockClientManager(mockCtrl)
					mockClientManager.EXPECT().IsClientInitialized().Return(true)
					mockClientManager.EXPECT().CloseClient().Return(errors.New("ClientManager.CloseClient() failed"))
					presenter.Print = func(writer io.Writer, output string) error {
						if writer == o.Stderr && strings.Contains(output, "ClientManager.CloseClient() failed") {
							return errors.New("Print() failed")
						}
						return origPrint(writer, output)
					}
					c := &cli{
						Exit:          exit,
						Cobra:         proxy.NewCobra(),
						RootCommand:   mockCommand,
						Context:       ctx,
						ClientManager: mockClientManager,
					}
					if got := c.Run(); got != 1 {
						t.Errorf("cli.Run() = %v, want %v", got, 1)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "",
			wantErr:    false,
			setup: func() {
				output = ""
			},
			cleanup: func() {
				output = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			output = ""
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			c := utility.NewCapturer(tt.fields.os, tt.fields.StdBuffer, tt.fields.ErrBuffer)
			gotStdOut, gotStdErr, err := c.CaptureOutput(func() {
				tt.args.fnc(mockCtrl)
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStdOut != tt.wantStdOut {
				t.Errorf("Run() gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("Run() gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}
