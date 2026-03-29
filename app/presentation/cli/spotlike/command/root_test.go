package command

import (
	"context"
	"testing"

	c "github.com/spf13/cobra"

	baseConfig "github.com/yanosea/spotlike/app/config"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command/spotlike"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/config"

	"github.com/yanosea/spotlike/pkg/proxy"
)

func TestNewRootCommand(t *testing.T) {
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
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				exit:    func(code int) {},
				cobra:   proxy.NewCobra(),
				version: "0.0.0",
				conf: &config.SpotlikeCliConfig{
					SpotlikeConfig: baseConfig.SpotlikeConfig{
						SpotifyID:           "test_id",
						SpotifySecret:       "test_secret",
						SpotifyRedirectUri:  "http://localhost:8080/callback",
						SpotifyRefreshToken: "test_refresh_token",
					},
				},
				output: &output,
			},
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
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			got := NewRootCommand(tt.args.exit, tt.args.cobra, tt.args.version, tt.args.conf, tt.args.output)
			if got == nil {
				t.Errorf("NewRootCommand() = %v, want not nil", got)
			} else {
				if err := got.RunE(nil, []string{}); err != nil {
					t.Errorf("Failed to run the root command: %v", err)
				}
			}
		})
	}
}

func Test_runRoot(t *testing.T) {
	type args struct {
		cmd        *c.Command
		cmdArgs    []string
		versionCmd proxy.Command
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		setup   func(tt *args)
		cleanup func()
	}{
		{
			name: "positive testing (rootOps.Version is true)",
			args: args{
				cmd:     &c.Command{},
				cmdArgs: []string{},
				versionCmd: spotlike.NewVersionCommand(
					proxy.NewCobra(),
					"0.0.0",
					&output,
				),
			},
			want:    "spotlike version 0.0.0",
			wantErr: false,
			setup: func(_ *args) {
				rootOps.Version = true
				output = ""
			},
			cleanup: func() {
				rootOps.Version = false
				output = ""
			},
		},
		{
			name: "positive testing (rootOps.Version is false)",
			args: args{
				cmd:     nil,
				cmdArgs: []string{},
				versionCmd: spotlike.NewVersionCommand(
					proxy.NewCobra(),
					"0.0.0",
					&output,
				),
			},
			want:    "",
			wantErr: false,
			setup: func(tt *args) {
				output = ""
				// Reset ClientManager first to avoid "client already initialized" error
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
				cm := api.NewClientManager(
					proxy.NewSpotify(),
					proxy.NewHttp(),
					proxy.NewRandstr(),
					proxy.NewUrl(),
				)
				if err := cm.InitializeClient(
					context.Background(),
					&api.ClientConfig{
						SpotifyID:           "test_id",
						SpotifySecret:       "test_secret",
						SpotifyRedirectUri:  "http://localhost:8080/callback",
						SpotifyRefreshToken: "test_refresh_token",
					},
				); err != nil {
					t.Errorf("Failed to initialize client: %v", err)
				}
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				tt.cmd = cmd
			},
			cleanup: func() {
				rootOps.Version = false
				output = ""
				if err := api.ResetClientManager(); err != nil {
					t.Errorf("Failed to reset client manager: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(&tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if err := runRoot(tt.args.cmd, tt.args.cmdArgs, &output, tt.args.versionCmd); (err != nil) != tt.wantErr {
				t.Errorf("runRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(output) == 0 {
				t.Errorf("runRoot() output is empty")
			}
			if tt.want != "" && output != tt.want {
				t.Errorf("runRoot() output = %v, want %v", output, tt.want)
			}
		})
	}
}
