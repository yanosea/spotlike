package unlike

import (
	"os"

	c "github.com/spf13/cobra"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/repository"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/presenter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// UnlikeArtistOptions represents the options for the unlike artist command.
type UnlikeArtistOptions struct {
	NoConfirm bool
	Format    string
}

var (
	// unlikeArtistOps is a variable to store the unlike artist options with the default values for injecting the dependencies in testing.
	unlikeArtistOps = UnlikeArtistOptions{
		NoConfirm: false,
		Format:    "table",
	}
)

// NewUnlikeArtistCommand creates a new unlike artist command.
func NewUnlikeArtistCommand(
	exit func(int),
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("artist")
	cmd.SetAliases([]string{"ar", "A"})
	cmd.SetUsageTemplate(unlikeArtistUsageTemplate)
	cmd.SetHelpTemplate(unlikeArtistHelpTemplate)
	cmd.SetSilenceErrors(true)
	cmd.Flags().BoolVarP(
		&unlikeArtistOps.NoConfirm,
		"no-confirm",
		"",
		false,
		"🚫 do not confirm before unliking the artist",
	)
	cmd.Flags().StringVarP(
		&unlikeArtistOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g: \"plain\")",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runUnlikeArtist(exit, cmd, authCmd, output, args)
		},
	)

	return cmd
}

// runUnlikeArtist executes the unlike artist command.
func runUnlikeArtist(exit func(int), cmd *c.Command, authCmd proxy.Command, output *string, args []string) error {
	if len(args) == 0 {
		o := formatter.Yellow("⚡ No ID arguments specified...")
		*output = o
		return nil
	}

	clientManager := api.GetClientManager()
	if clientManager == nil {
		o := formatter.Red("❌ Client manager is not initialized...")
		*output = o
		return nil
	}
	_, err := clientManager.GetClient()
	if err != nil && err.Error() == "client not initialized" {
		if err := authCmd.RunE(cmd, args); err != nil {
			return err
		}
	} else if err != nil {
		o := formatter.Red("❌ Failed to get client...")
		*output = o
		return err
	}

	var gAucoDtos []*spotlikeApp.GetArtistUseCaseOutputDto
	artistRepo := repository.NewArtistRepository()
	gAuc := spotlikeApp.NewGetArtistUseCase(artistRepo)
	for _, id := range args {
		gAucoDto, err := gAuc.Run(cmd.Context(), id)
		if err != nil && err.Error() != "Resource not found" {
			return err
		}

		if gAucoDto == nil {
			o := formatter.Yellow("⚡ The id " + id + " is not found or it is not an artist...")
			*output = o
			continue
		}

		gAucoDtos = append(gAucoDtos, gAucoDto)
	}

	clAuc := spotlikeApp.NewCheckLikeArtistUseCase(artistRepo)
	uAuc := spotlikeApp.NewUnlikeArtistUseCase(artistRepo)
	var unlikeExecutedArtists []*spotlikeApp.GetArtistUseCaseOutputDto
	for _, gAucoDto := range gAucoDtos {
		alreadyLiked, err := clAuc.Run(cmd.Context(), gAucoDto.ID)
		if err != nil {
			return err
		}
		if !alreadyLiked {
			if err := presenter.Print(os.Stdout, formatter.Blue("⏩ Artist "+gAucoDto.Name+" ("+gAucoDto.ID+") "+"is not liked. skipping...")); err != nil {
				return err
			}
			continue
		}

		if !unlikeArtistOps.NoConfirm {
			if answer, err := presenter.RunPrompt(
				"Proceed with unliking " + gAucoDto.Name + " (" + gAucoDto.ID + ") ? [y/N]",
			); err != nil && err.Error() == "^C" {
				if err := presenter.Print(os.Stdout, "\n"); err != nil {
					return err
				}
				if err := presenter.Print(os.Stdout, formatter.Yellow("🚫 Cancelled unliking...")); err != nil {
					return err
				}
				exit(130)
				return nil
			} else if err != nil {
				return err
			} else if answer != "y" && answer != "Y" {
				o := formatter.Yellow("🚫 Cancelled unliking artist " + gAucoDto.Name + " (" + gAucoDto.ID + ") ...")
				*output = o
				continue
			}
		}

		if err := uAuc.Run(cmd.Context(), gAucoDto.ID); err != nil {
			return err
		}

		unlikeExecutedArtists = append(unlikeExecutedArtists, gAucoDto)
	}

	if len(unlikeExecutedArtists) != 0 {
		f, err := formatter.NewFormatter(unlikeArtistOps.Format)
		if err != nil {
			o := formatter.Red("❌ Failed to create a formatter...")
			*output = o
			return err
		}
		o, err := f.Format(unlikeExecutedArtists)
		if err != nil {
			return err
		}
		*output = "\n" + o
		if err := presenter.Print(os.Stdout, formatter.Green("✅💔🎤 Successfully unliked artists below!")); err != nil {
			return err
		}
	}

	return nil
}

const (
	// unlikeArtistHelpTemplate is the help template of the unlike artist command.
	unlikeArtistHelpTemplate = `💔🎤 Unlike artists on Spotify by ID.

You can unlike artists on Spotify by ID.

Before using this command,
you need to get the ID of the artist you want to like by using the search command.

` + unlikeArtistUsageTemplate
	// unlikeArtistUsageTemplate is the usage template of the unlike artist command.
	unlikeArtistUsageTemplate = `Usage:
  spotlike unlike artist [flags] [arguments]
  spotlike unlike ar     [flags] [arguments]
  spotlike unlike A      [flags] [arguments]

Flags:
  --no-confirm  🚫 do not confirm before unliking the artist
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for artist

Arguments:
  ID  🆔 ID of the artists (e.g. : "00DuPiLri3mNomvvM3nZvU 3B9O5mYYw89fFXkwKh7jCS")
`
)
