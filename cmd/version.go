package cmd

import (
	"fmt"
	"github.com/yanosea/spotlike/util"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	version_use              = "version"
	version_short            = "Show the version of spotlike."
	version_message_template = "spotlike version %s"
)

func newVersionCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   version_use,
		Short: version_short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return globalOption.version()
		},
	}

	cmd.SetOut(globalOption.Out)

	return cmd
}

func (g *GlobalOption) version() error {
	util.PrintlnWithWriter(g.Out, fmt.Sprintf(version_message_template, version))
	return nil
}
