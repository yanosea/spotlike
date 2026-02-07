package like

import (
	"fmt"
	"os"

	c "github.com/spf13/cobra"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/repository"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/presenter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// LikeTrackOptions represents the options for the like track command.
type LikeTrackOptions struct {
	Artist    string
	Album     string
	NoConfirm bool
	Format    string
}

var (
	// likeTrackOps is a variable to store the like track options with the default values for injecting the dependencies in testing.
	likeTrackOps = LikeTrackOptions{
		Artist:    "",
		Album:     "",
		NoConfirm: false,
		Format:    "table",
	}
)

// NewLikeTrackCommand creates a new like track command.
func NewLikeTrackCommand(
	exit func(int),
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("track")
	cmd.SetAliases([]string{"tr", "t"})
	cmd.SetUsageTemplate(likeTrackUsageTemplate)
	cmd.SetHelpTemplate(likeTrackHelpTemplate)
	cmd.SetSilenceErrors(true)
	cmd.Flags().StringVarP(
		&likeTrackOps.Artist,
		"artist",
		"A",
		"",
		"🆔 an ID of the artist to like all albums released by the artist",
	)
	cmd.Flags().StringVarP(
		&likeTrackOps.Album,
		"album",
		"a",
		"",
		"🆔 an ID of the album to like all tracks in the album",
	)
	cmd.Flags().BoolVarP(
		&likeTrackOps.NoConfirm,
		"no-confirm",
		"",
		false,
		"🚫 do not confirm before liking the track",
	)
	cmd.Flags().StringVarP(
		&likeTrackOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g: \"plain\")",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runLikeTrack(exit, cmd, authCmd, output, args)
		},
	)

	return cmd
}

// runLikeTrack executes the like track command.
func runLikeTrack(exit func(int), cmd *c.Command, authCmd proxy.Command, output *string, args []string) error {
	if likeTrackOps.Artist != "" && likeTrackOps.Album != "" {
		o := formatter.Yellow("⚡ Both artist and album flags can not be specified at the same time...")
		*output = o
		return nil
	}

	if likeTrackOps.Artist == "" && likeTrackOps.Album == "" && len(args) == 0 {
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

	var gtucoDtos []*spotlikeApp.GetTrackUseCaseOutputDto
	trackRepo := repository.NewTrackRepository()
	if likeTrackOps.Artist != "" {
		artistRepo := repository.NewArtistRepository()
		gAuc := spotlikeApp.NewGetArtistUseCase(artistRepo)
		gAucoDto, err := gAuc.Run(cmd.Context(), likeTrackOps.Artist)
		if err != nil && err.Error() != "Resource not found" {
			return err
		}

		if gAucoDto == nil {
			o := formatter.Yellow("⚡ The id " + likeTrackOps.Artist + " is not found or it is not an artist...")
			*output = o
			return nil
		}

		gatAuc := spotlikeApp.NewGetAllTracksByArtistIdUseCase(trackRepo)
		gatAucoDtos, err := gatAuc.Run(cmd.Context(), likeTrackOps.Artist)
		if err != nil {
			return err
		}

		for _, gatAucoDto := range gatAucoDtos {
			gtucoDtos = append(
				gtucoDtos,
				&spotlikeApp.GetTrackUseCaseOutputDto{
					ID:          gatAucoDto.ID,
					Name:        gatAucoDto.Name,
					Artists:     gatAucoDto.Artists,
					Album:       gatAucoDto.Album,
					TrackNumber: gatAucoDto.TrackNumber,
					ReleaseDate: gatAucoDto.ReleaseDate,
				},
			)
		}
	} else if likeTrackOps.Album != "" {
		albumRepo := repository.NewAlbumRepository()
		gauc := spotlikeApp.NewGetAlbumUseCase(albumRepo)
		gaucoDto, err := gauc.Run(cmd.Context(), likeTrackOps.Album)
		if err != nil && err.Error() != "Resource not found" {
			return err
		}

		if gaucoDto == nil {
			o := formatter.Yellow("⚡ The id " + likeTrackOps.Album + " is not found or it is not an album...")
			*output = o
			return nil
		}

		gatauc := spotlikeApp.NewGetAllTracksByAlbumIdUseCase(trackRepo)
		gataucoDtos, err := gatauc.Run(cmd.Context(), likeTrackOps.Album)
		if err != nil {
			return err
		}

		for _, gataucoDto := range gataucoDtos {
			gtucoDtos = append(
				gtucoDtos,
				&spotlikeApp.GetTrackUseCaseOutputDto{
					ID:          gataucoDto.ID,
					Name:        gataucoDto.Name,
					Artists:     gataucoDto.Artists,
					Album:       gataucoDto.Album,
					TrackNumber: gataucoDto.TrackNumber,
					ReleaseDate: gataucoDto.ReleaseDate,
				},
			)
		}
	} else {
		for _, id := range args {
			gtuc := spotlikeApp.NewGetTrackUseCase(trackRepo)
			gtucoDto, err := gtuc.Run(cmd.Context(), id)
			if err != nil && err.Error() != "Resource not found" {
				return err
			}

			if gtucoDto == nil {
				o := formatter.Yellow("⚡ The id " + id + " is not found or it is not a track...")
				*output = o
				continue
			}

			gtucoDtos = append(gtucoDtos, gtucoDto)
		}
	}

	cltuc := spotlikeApp.NewCheckLikeTrackUseCase(trackRepo)
	ltuc := spotlikeApp.NewLikeTrackUseCase(trackRepo)
	var likeExecutedTracks []*spotlikeApp.GetTrackUseCaseOutputDto
	for _, gtucoDto := range gtucoDtos {
		alreadyLiked, err := cltuc.Run(cmd.Context(), gtucoDto.ID)
		if err != nil {
			return err
		}
		if alreadyLiked {
			if err := presenter.Print(os.Stdout, formatter.Blue("⏩ Track #"+fmt.Sprint(gtucoDto.TrackNumber)+" "+gtucoDto.Name+" ("+gtucoDto.ID+")"+" on "+gtucoDto.Album+" rereased by "+gtucoDto.Artists+" is already liked. skipping...")); err != nil {
				return err
			}
			continue
		}

		if !likeTrackOps.NoConfirm {
			if answer, err := presenter.RunPrompt(
				"Proceed with liking " + gtucoDto.Name + " (" + gtucoDto.ID + ") ? [y/N]",
			); err != nil && err.Error() == "^C" {
				if err := presenter.Print(os.Stdout, "\n"); err != nil {
					return err
				}
				if err := presenter.Print(os.Stdout, formatter.Yellow("🚫 Cancelled liking...")); err != nil {
					return err
				}
				exit(130)
				return nil
			} else if err != nil {
				return err
			} else if answer != "y" && answer != "Y" {
				o := formatter.Yellow("🚫 Cancelled liking track " + gtucoDto.Name + " (" + gtucoDto.ID + ") ...")
				*output = o
				continue
			}
		}

		if err := ltuc.Run(cmd.Context(), gtucoDto.ID); err != nil {
			return err
		}

		likeExecutedTracks = append(likeExecutedTracks, gtucoDto)
	}

	if len(likeExecutedTracks) != 0 {
		f, err := formatter.NewFormatter(likeTrackOps.Format)
		if err != nil {
			o := formatter.Red("❌ Failed to create a formatter...")
			*output = o
			return err
		}
		o, err := f.Format(likeExecutedTracks)
		if err != nil {
			return err
		}
		*output = "\n" + o
		if err := presenter.Print(os.Stdout, formatter.Green("✅🤍🎵 Successfully liked tracks below!")); err != nil {
			return err
		}
	}

	return nil
}

const (
	// likeTrackHelpTemplate is a template for the help message of the like track command.
	likeTrackHelpTemplate = `🤍🎵 Like tracks on Spotify by ID.

You can like tracks on Spotify by ID.

Before using this command,
you need to get the ID of the album you want to like by using the search command.

Also, you can like all tracks released by the artist with specifying the ID of the artist with artist flag.
If you specify artist flag, the arguments would be ignored.
Also, you can like all tracks in the album with specifying the ID of the album with album flag.
If you specify album flag, the arguments would be ignored.
Both artist and album flags can not be specified at the same time.

` + likeTrackUsageTemplate
	// likeTrackUsageTemplate is a template for the usage message of the like track command.
	likeTrackUsageTemplate = `Usage:
  spotlike like track [flags] [arguments]
  spotlike like tr    [flags] [arguments]
  spotlike like t     [flags] [arguments]

Flags:
  -A, --artist  🆔 an ID of the artist to like all albums released by the artist
  -a, --album   🆔 an ID of the album to like all tracks in the album
  --no-confirm  🚫 do not confirm before liking the track
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for track

Arguments:
  ID  🆔 ID of the tracks (e.g: "20q73dOrP7ceLGAJQVtuTq 5A7nqzXUt5IZIOA7oNBv6M")
`
)
