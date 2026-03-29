package like

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

func TestNewLikeCommand(t *testing.T) {
	output := ""
	exit := os.Exit

	type args struct {
		exit    func(int)
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
				exit:  exit,
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
			got := NewLikeCommand(tt.args.exit, tt.args.cobra, tt.args.authCmd, tt.args.output)
			if got == nil {
				t.Errorf("NewLikeCommand() = %v, want %v", got, tt.want)
			} else {
				cmd := &c.Command{}
				cmd.SetContext(context.Background())
				if err := got.RunE(cmd, []string{}); err != nil {
					t.Errorf("Failed to run like command: %v", err)
				}
			}
		})
	}
}

func Test_runLike(t *testing.T) {
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

  - 🎤 artist
  - 💿 album
  - 🎵 track

Use "spotlike like --help" for more information about spotlike like.
Use "spotlike like [command] --help" for more information about a command.
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runLike(tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runLike() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *tt.args.output != tt.wantOutput {
				t.Errorf("runLike() output = %v, want %v", *tt.args.output, tt.wantOutput)
			}
		})
	}
}
