package get

import (
	c "github.com/spf13/cobra"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/repository"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// GetTracksOptions represents the options for the get tracks command.
type GetTracksOptions struct {
	Format string
}

var (
	// getTracksOps is a variable to store the get tracks options with the default values for injecting the dependencies in testing.
	getTracksOps = GetTracksOptions{
		Format: "table",
	}
)

// NewGetTracksCommand returns a new instance of the get tracks command.
func NewGetTracksCommand(
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("tracks")
	cmd.SetAliases([]string{"trs", "t"})
	cmd.SetUsageTemplate(getTracksUsageTemplate)
	cmd.SetHelpTemplate(getTracksHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.Flags().StringVarP(
		&getTracksOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g: \"plain\")",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runGetTracks(cmd, authCmd, output, args)
		},
	)

	return cmd
}

// runGetTracks is a function to get the tracks on Spotify by ID.
func runGetTracks(cmd *c.Command, authCmd proxy.Command, output *string, args []string) error {
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

	albumRepo := repository.NewAlbumRepository()
	gauc := spotlikeApp.NewGetAlbumUseCase(albumRepo)
	gaucoDto, err := gauc.Run(cmd.Context(), args[0])
	if err != nil && err.Error() != "Resource not found" {
		return err
	}

	if gAucoDto == nil && gaucoDto == nil {
		o := formatter.Yellow("⚡ The id " + args[0] + " is not found... or it is not an artist or album...")
		*output = o
		return nil
	}

	var tracks []*spotlikeApp.GetTrackUseCaseOutputDto
	trackRepo := repository.NewTrackRepository()
	if gAucoDto != nil {
		gatAuc := spotlikeApp.NewGetAllTracksByArtistIdUseCase(trackRepo)
		gatAucoDtos, err := gatAuc.Run(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		for _, gatAucoDto := range gatAucoDtos {
			track := &spotlikeApp.GetTrackUseCaseOutputDto{
				ID:          gatAucoDto.ID,
				Name:        gatAucoDto.Name,
				Artists:     gatAucoDto.Artists,
				Album:       gatAucoDto.Album,
				TrackNumber: gatAucoDto.TrackNumber,
				ReleaseDate: gatAucoDto.ReleaseDate,
			}
			tracks = append(tracks, track)
		}
	}
	if gaucoDto != nil {
		gatAuc := spotlikeApp.NewGetAllTracksByAlbumIdUseCase(trackRepo)
		gatAucoDtos, err := gatAuc.Run(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		for _, gatAucoDto := range gatAucoDtos {
			track := &spotlikeApp.GetTrackUseCaseOutputDto{
				ID:          gatAucoDto.ID,
				Name:        gatAucoDto.Name,
				Artists:     gatAucoDto.Artists,
				Album:       gatAucoDto.Album,
				TrackNumber: gatAucoDto.TrackNumber,
				ReleaseDate: gatAucoDto.ReleaseDate,
			}
			tracks = append(tracks, track)
		}
	}

	f, err := formatter.NewFormatter(getTracksOps.Format)
	if err != nil {
		o := formatter.Red("❌ Failed to create a formatter...")
		*output = o
		return err
	}
	o, err := f.Format(tracks)
	if err != nil {
		return err
	}
	*output = "\n" + o

	return nil
}

const (
	// getTracksHelpTemplate is a help template for the get tracks command.
	getTracksHelpTemplate = `ℹ️🎵 Get the information of the tracks on Spotify by ID.

You can get the information of the tracks on Spotify by artist ID or album ID.

Before using this command,
you need to get the ID of the artist or album you want to get by using the search command.

` + getTracksUsageTemplate
	// getTracksUsageTemplate is a usage template for the get tracks command.
	getTracksUsageTemplate = `Usage:
  spotlike get tracks [flags] [arguments]
  spotlike get trs    [flags] [arguments]
  spotlike get t      [flags] [arguments]

Flags:
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for tracks

Argument:
  ID  🆔 ID of the artist or album (e.g: "00DuPiLri3mNomvvM3nZvU")
`
)
