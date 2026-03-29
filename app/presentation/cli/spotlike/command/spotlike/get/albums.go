package get

import (
	c "github.com/spf13/cobra"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/repository"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// GetAlbumsOptions represents the options for the get albums command.
type GetAlbumsOptions struct {
	Format string
}

var (
	// getAlbumsOps is a variable to store the get albums options with the default values for injecting the dependencies in testing.
	getAlbumsOps = GetAlbumsOptions{
		Format: "table",
	}
)

// NewGetAlbumsCommand returns a new instance of the get albums command.
func NewGetAlbumsCommand(
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("albums")
	cmd.SetAliases([]string{"als", "a"})
	cmd.SetUsageTemplate(getAlbumsUsageTemplate)
	cmd.SetHelpTemplate(getAlbumsHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.Flags().StringVarP(
		&getAlbumsOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g: \"plain\")",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runGetAlbums(cmd, authCmd, output, args)
		},
	)

	return cmd
}

// runGetAlbums is a function to get the albums on Spotify by artist ID.
func runGetAlbums(cmd *c.Command, authCmd proxy.Command, output *string, args []string) error {
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
	}

	artistRepo := repository.NewArtistRepository()
	gAuc := spotlikeApp.NewGetArtistUseCase(artistRepo)
	gAucoDto, err := gAuc.Run(cmd.Context(), args[0])
	if err != nil && err.Error() != "Resource not found" {
		return err
	}

	if gAucoDto == nil {
		o := formatter.Yellow("⚡ The id " + args[0] + " is not found... or it is not an artist...")
		*output = o
		return nil
	}

	albumRepo := repository.NewAlbumRepository()
	gaauc := spotlikeApp.NewGetAllAlbumsByArtistIdUseCase(albumRepo)
	gaaucoDtos, err := gaauc.Run(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	var albums []*spotlikeApp.GetAlbumUseCaseOutputDto
	for _, gaaucoDto := range gaaucoDtos {
		album := &spotlikeApp.GetAlbumUseCaseOutputDto{
			ID:          gaaucoDto.ID,
			Name:        gaaucoDto.Name,
			Artists:     gaaucoDto.Artists,
			ReleaseDate: gaaucoDto.ReleaseDate,
		}
		albums = append(albums, album)
	}

	f, err := formatter.NewFormatter(getAlbumsOps.Format)
	if err != nil {
		o := formatter.Red("❌ Failed to create a formatter...")
		*output = o
		return err
	}
	o, err := f.Format(albums)
	if err != nil {
		return err
	}
	*output = "\n" + o

	return nil
}

const (
	// getAlbumsHelpTemplate is a help template for the get albums command.
	getAlbumsHelpTemplate = `ℹ️💿 Get the information of the albums on Spotify by artist ID.

You can get the information of the albums on Spotify by artist ID.

Before using this command,
you need to get the ID of the artist you want to get by using the search command.

` + getAlbumsUsageTemplate
	// getAlbumsUsageTemplate is a usage template for the get albums command.
	getAlbumsUsageTemplate = `Usage:
  spotlike get albums [flags] [arguments]
  spotlike get als    [flags] [arguments]
  spotlike get a      [flags] [arguments]

Flags:
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for albums

Argument:
  ID  🆔 ID of the artist (e.g. : "00DuPiLri3mNomvvM3nZvU")
`
)
