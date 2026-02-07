package main

import (
	"context"
	o "os"
	"testing"

	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/command"

	"github.com/yanosea/spotlike/pkg/proxy"
	"github.com/yanosea/spotlike/pkg/utility"

	"go.uber.org/mock/gomock"
)

func Test_main(t *testing.T) {
	os := proxy.NewOs()
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()
	if err := o.Setenv("SPOTIFY_ID", "test_id"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := o.Setenv("SPOTIFY_SECRET", "test_secret"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := o.Setenv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/callback"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	if err := o.Setenv("SPOTIFY_REFRESH_TOKEN", "test_refresh_token"); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	origExit := exit
	exit = func(code int) {}
	defer func() {
		exit = origExit
	}()
	origArgs := o.Args

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
		wantAny    bool
		wantStdOut string
		wantStdErr string
		wantErr    bool
		setup      func(mockCtrl *gomock.Controller)
		cleanup    func()
	}{
		{
			name: "initial execution (spotlike)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/spotlike"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: "⚡ Use sub command below...\n\n" +
				"- 🔑 auth,       au,   a - Authenticate your Spotify client.\n" +
				"- 📚 get,        ge,   g - Get the information of the content on Spotify by ID.\n" +
				"- 🤍 like,       li,   l - Like content on Spotify by ID.\n" +
				"- 💔 unlike,     un,   u - Unlike content on Spotify by ID.\n" +
				"- 🔍 search,     se,   s - Search for the ID of content in Spotify.\n" +
				"- 🔧 completion, comp, c - Generate the autocompletion script for the specified shell.\n" +
				"- 🔖 version,    ver,  v - Show the version of spotlike.\n" +
				"- 🤝 help                - Help for spotlike.\n\n" +
				"Use \"spotlike --help\" for more information about spotlike.\n" +
				"Use \"spotlike [command] --help\" for more information about a command.\n\n",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "version execution (spotlike --version)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/spotlike", "--version"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantStdOut: "spotlike version (devel)\n",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "help execution (spotlike --help)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/spotlike", "--help"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    true,
			wantStdOut: "",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "auth execution (spotlike auth)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/spotlike", "auth"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    true,
			wantStdOut: "",
			wantStdErr: "",
			setup:      nil,
			cleanup:    nil,
		},
		{
			name: "negative testing (cli.Init() failed)",
			fields: fields{
				Os:        os,
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					o.Args = []string{"/path/to/spotlike"}
					defer func() {
						o.Args = origArgs
					}()
					main()
				},
			},
			wantAny:    false,
			wantStdOut: "",
			wantStdErr: "",
			setup: func(mockCtrl *gomock.Controller) {
				mockCli := command.NewMockCli(mockCtrl)
				mockCli.EXPECT().Init(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1)
				mockCli.EXPECT().Run()
				origNewCli := command.NewCli
				command.NewCli = func(exit func(int), cobra proxy.Cobra, ctx context.Context) command.Cli {
					return mockCli
				}
				t.Cleanup(func() {
					command.NewCli = origNewCli
				})
			},
			cleanup: func() {
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
			if !tt.wantAny {
				if gotStdOut != tt.wantStdOut {
					t.Errorf("main() gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
				}
				if gotStdErr != tt.wantStdErr {
					t.Errorf("main() gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
				}
			} else {
				t.Logf("main() gotStdOut = %v", gotStdOut)
				t.Logf("main() gotStdErr = %v", gotStdErr)
			}
		})
	}
}

