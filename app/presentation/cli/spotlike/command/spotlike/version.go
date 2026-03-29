package spotlike

import (
	c "github.com/spf13/cobra"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

var (
	// format is the format of the output.
	format = "plain"
)

// NewVersionCommand returns a new instance of the version command.
func NewVersionCommand(
	cobra proxy.Cobra,
	version string,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("version")
	cmd.SetAliases([]string{"ver", "v"})
	cmd.SetUsageTemplate(versionUsageTemplate)
	cmd.SetHelpTemplate(versionHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(
		func(_ *c.Command, _ []string) error {
			return runVersion(version, output)
		},
	)

	return cmd
}

// runVersion runs the version command.
func runVersion(version string, output *string) error {
	uc := spotlikeApp.NewGetVersionUseCase()
	dto := uc.Run(version)

	f, err := formatter.NewFormatter(format)
	if err != nil {
		o := formatter.Red("❌ Failed to create a formatter...")
		*output = o
		return err
	}
	o, err := f.Format(dto)
	if err != nil {
		return err
	}
	*output = o

	return nil
}

const (
	// versionHelpTemplate is the help template of the version command.
	versionHelpTemplate = `🔖 Show the version of spotlike

` + versionUsageTemplate
	// versionUsageTemplate is the usage template of the version command.
	versionUsageTemplate = `Usage:
  spotlike version [flags]
  spotlike ver     [flags]
  spotlike v       [flags]

Flags:
  -h, --help  🤝 help for version
`
)
