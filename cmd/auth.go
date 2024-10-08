package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/library/auth"
	"github.com/yanosea/spotlike/app/library/utility"
	"github.com/yanosea/spotlike/app/proxy/cobra"
	"github.com/yanosea/spotlike/app/proxy/color"
	"github.com/yanosea/spotlike/app/proxy/io"
	"github.com/yanosea/spotlike/cmd/constant"
)

type authOption struct {
	Out     ioproxy.WriterInstanceInterface
	ErrOut  ioproxy.WriterInstanceInterface
	Args    []string
	Utility utility.UtilityInterface
}

func NewAuthCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &authOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    g.Args,
		Utility: g.Utility,
	}

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.AUTH_USE
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.authRunE

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.AUTH_HELP_TEMPLATE)

	return cmd
}

func (o *authOption) authRunE(_ *cobra.Command, _ []string) error {
	// check if auth info is already set
	if auth.IsEnvsSet() {
		colorProxy := colorproxy.New()
		// if already set, print message and return
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.AUTH_MESSAGE_ALREADY_AUTHENTICATED))
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
