package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/library/utility"
	"github.com/yanosea/spotlike/app/proxy/cobra"
	"github.com/yanosea/spotlike/app/proxy/color"
	"github.com/yanosea/spotlike/app/proxy/io"
	"github.com/yanosea/spotlike/cmd/constant"
)

type likeOption struct {
	Out     ioproxy.WriterInstanceInterface
	ErrOut  ioproxy.WriterInstanceInterface
	Args    []string
	Utility utility.UtilityInterface
}

func NewLikeCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &likeOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    g.Args,
		Utility: g.Utility,
	}

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.LIKE_USE
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.likeRunE

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.LIKE_HELP_TEMPLATE)

	cmd.AddCommand(
		NewLikeArtistCommand(g),
		NewLikeAlbumCommand(g),
		NewLikeTrackCommand(g),
	)

	return cmd
}

func (o *likeOption) likeRunE(_ *cobra.Command, _ []string) error {
	colorproxy := colorproxy.New()
	// if no sub command is specified, print the message and return nil.
	o.Utility.PrintlnWithWriter(o.Out, colorproxy.YellowString(constant.LIKE_MESSAGE_NO_SUB_COMMAND))
	return nil
}
