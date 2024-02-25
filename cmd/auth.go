package cmd

import (
	"io"

	"github.com/yanosea/spotlike/auth"
	"github.com/yanosea/spotlike/util"

	// https://github.com/fatih/color
	"github.com/fatih/color"
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	auth_help_template = `🔑 Authenticate your Spotify client.

You have to authenticate your Spotify client to use spotlike at first.
spotlike will ask you to input your Client ID, Client Secret, Redirect URI, and Refresh Token.

Usage:
  spotlike auth [flags]

Flags:
  -h, --help   help for auth
`
	auth_use   = "auth"
	auth_short = "🔑 Authenticate your Spotify client."
	auth_long  = `🔑 Authenticate your Spotify client.

You have to authenticate your Spotify client to use spotlike at first.
spotlike will ask you to input your Client ID, Client Secret, Redirect URI, and Refresh Token.`
	auth_message_already_authenticated = "✅ You are already authenticated and set envs!"
)

type authOption struct {
	Out    io.Writer
	ErrOut io.Writer
}

func newAuthCommand(globalOption *GlobalOption) *cobra.Command {
	o := &authOption{}
	cmd := &cobra.Command{
		Use:   auth_use,
		Short: auth_short,
		Long:  auth_long,
		RunE: func(cmd *cobra.Command, args []string) error {

			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut

			return o.auth()
		},
	}

	o.Out = globalOption.Out
	o.ErrOut = globalOption.ErrOut
	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	cmd.SetHelpTemplate(auth_help_template)

	return cmd
}

func (o *authOption) auth() error {
	// check if auth info is already set
	if auth.IsEnvsSet() {
		// if already set, print message and return
		util.PrintlnWithWriter(o.Out, color.YellowString(auth_message_already_authenticated))
		return nil
	}

	// set auth info
	auth.SetAuthInfo()

	// execute auth
	_, err := auth.Authenticate(o.Out)
	if err != nil {
		return err
	}

	return nil
}
