package get

import (
	"context"
	"os"
	"testing"

	c "github.com/spf13/cobra"

	baseconfig "github.com/yanosea/spotlike/app/config"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command/spotlike"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/config"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

func TestNewGetCommand(t *testing.T) {
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
			got := NewGetCommand(tt.args.cobra, tt.args.authCmd, tt.args.output)
			if got == nil {
				t.Errorf("NewGetCommand() = %v, want %v", got, tt.want)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run get command: %v", err)
				}
			}
		})
	}
}

func Test_runGet(t *testing.T) {
	output := ""

	type args struct {
		output *string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
	}{
		{
			name: "positive testing",
			args: args{
				output: &output,
			},
			wantOutput: formatter.Yellow("⚡ Use sub command below...") + `

  - 💿 album
  - 🎵 track

Use "spotlike get --help" for more information about spotlike get.
Use "spotlike get [command] --help" for more information about a command.
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runGet(tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runGet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *tt.args.output != tt.wantOutput {
				t.Errorf("runGet() output = %v, want %v", *tt.args.output, tt.wantOutput)
			}
		})
	}
}
