package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/library/utility"
	"github.com/yanosea/spotlike/app/library/versionprovider"
	"github.com/yanosea/spotlike/app/proxy/cobra"
	"github.com/yanosea/spotlike/app/proxy/color"
	"github.com/yanosea/spotlike/app/proxy/debug"
	"github.com/yanosea/spotlike/app/proxy/fmt"
	"github.com/yanosea/spotlike/app/proxy/io"
	"github.com/yanosea/spotlike/app/proxy/os"
	"github.com/yanosea/spotlike/cmd/constant"
)

// ver is the version of the spotlike.
var ver = "develop"

type GlobalOption struct {
	Out            ioproxy.WriterInstanceInterface
	ErrOut         ioproxy.WriterInstanceInterface
	Args           []string
	Utility        utility.UtilityInterface
	NewRootCommand func(ow, ew ioproxy.WriterInstanceInterface, args []string) cobraproxy.CommandInstanceInterface
}

// rootOption is the struct for root command.
type rootOption struct {
	Out     ioproxy.WriterInstanceInterface
	ErrOut  ioproxy.WriterInstanceInterface
	Args    []string
	Utility utility.UtilityInterface
}

// NewGlobalOption creates a new global option.
func NewGlobalOption(fmtProxy fmtproxy.Fmt, osProxy osproxy.Os) *GlobalOption {
	return &GlobalOption{
		Out:    osproxy.Stdout,
		ErrOut: osproxy.Stderr,
		Args:   osproxy.Args[1:],
		Utility: utility.New(
			fmtProxy,
			osProxy,
		),
		NewRootCommand: NewRootCommand,
	}
}

// Execute executes the spotlike.
func (g *GlobalOption) Execute() int {
	rootCmd := g.NewRootCommand(g.Out, g.ErrOut, g.Args)
	if err := rootCmd.Execute(); err != nil {
		colorProxy := colorproxy.New()
		g.Utility.PrintlnWithWriter(g.ErrOut, colorProxy.RedString(err.Error()))
		return 1
	}
	return 0
}

// NewRootCommand creates a new root command.
func NewRootCommand(ow, ew ioproxy.WriterInstanceInterface, cmdArgs []string) cobraproxy.CommandInstanceInterface {
	g := &GlobalOption{
		Out:    ow,
		ErrOut: ew,
		Args:   cmdArgs,
	}
	o := &rootOption{
		Out:    g.Out,
		ErrOut: g.ErrOut,
		Args:   cmdArgs,
	}

	v := versionprovider.New(debugproxy.New())

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.ROOT_USE
	cmd.FieldCommand.Version = v.GetVersion(ver)
	cmd.FieldCommand.SilenceErrors = true
	cmd.FieldCommand.SilenceUsage = true
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.rootRunE

	cmd.SetOut(ow)
	cmd.SetErr(ew)
	cmd.SetHelpTemplate(constant.ROOT_HELP_TEMPLATE)

	cmd.AddCommand(
		NewAuthCommand(g),
		NewCompletionCommand(g),
		NewLikeCommand(g),
		NewSearchCommand(g),
		NewVersionCommand(g),
	)

	cmd.SetArgs(cmdArgs)

	return cmd
}

// rootRunE is the function to run root command.
func (o *rootOption) rootRunE(_ *cobra.Command, _ []string) error {
	colorproxy := colorproxy.New()
	// if no sub command is specified, print the message and return nil.
	o.Utility.PrintlnWithWriter(o.Out, colorproxy.YellowString(constant.ROOT_MESSAGE_NO_SUB_COMMAND))
	return nil
}
