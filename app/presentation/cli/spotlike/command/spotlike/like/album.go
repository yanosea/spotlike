package like

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

// LikeAlbumOptions represents the options for the like album command.
type LikeAlbumOptions struct {
	Artist    string
	NoConfirm bool
	Format    string
}

var (
	// likeAlbumOps is a variable to store the like album options with the default values for injecting the dependencies in testing.
	likeAlbumOps = LikeAlbumOptions{
		Artist:    "",
		NoConfirm: false,
		Format:    "table",
	}
)

// NewLikeAlbumCommand creates a new like album command.
func NewLikeAlbumCommand(
	exit func(int),
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("album")
	cmd.SetAliases([]string{"al", "a"})
	cmd.SetUsageTemplate(likeAlbumUsageTemplate)
	cmd.SetHelpTemplate(likeAlbumHelpTemplate)
	cmd.SetSilenceErrors(true)
	cmd.Flags().StringVarP(
		&likeAlbumOps.Artist,
		"artist",
		"A",
		"",
		"🆔 an ID of the artist to like all albums released by the artist",
	)
	cmd.Flags().BoolVarP(
		&likeAlbumOps.NoConfirm,
		"no-confirm",
		"",
		false,
		"🚫 do not confirm before liking the album",
	)
	cmd.Flags().StringVarP(
		&likeAlbumOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g: \"plain\")",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runLikeAlbum(exit, cmd, authCmd, output, args)
		},
	)

	return cmd
}

// runLikeAlbum runs the like album command.
func runLikeAlbum(exit func(int), cmd *c.Command, authCmd proxy.Command, output *string, args []string) error {
	if likeAlbumOps.Artist == "" && len(args) == 0 {
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

	var gaucoDtos []*spotlikeApp.GetAlbumUseCaseOutputDto
	albumRepo := repository.NewAlbumRepository()
	if likeAlbumOps.Artist != "" {
		artistRepo := repository.NewArtistRepository()
		gAuc := spotlikeApp.NewGetArtistUseCase(artistRepo)
		gAucoDto, err := gAuc.Run(cmd.Context(), likeAlbumOps.Artist)
		if err != nil && err.Error() != "Resource not found" {
			return err
		}

		if gAucoDto == nil {
			o := formatter.Yellow("⚡ The id " + likeAlbumOps.Artist + " is not found or it is not an artist...")
			*output = o
			return nil
		}

		gaaAuc := spotlikeApp.NewGetAllAlbumsByArtistIdUseCase(albumRepo)
		gaaAucoDtos, err := gaaAuc.Run(cmd.Context(), likeAlbumOps.Artist)
		if err != nil {
			return err
		}

		for _, gaaAucoDto := range gaaAucoDtos {
			gaucoDtos = append(
				gaucoDtos,
				&spotlikeApp.GetAlbumUseCaseOutputDto{
					ID:          gaaAucoDto.ID,
					Name:        gaaAucoDto.Name,
					Artists:     gaaAucoDto.Artists,
					ReleaseDate: gaaAucoDto.ReleaseDate,
				},
			)
		}
	} else {
		for _, id := range args {
			gauc := spotlikeApp.NewGetAlbumUseCase(albumRepo)
			gaucoDto, err := gauc.Run(cmd.Context(), id)
			if err != nil && err.Error() != "Resource not found" {
				return err
			}

			if gaucoDto == nil {
				o := formatter.Yellow("⚡ The id " + id + " is not found or it is not an album...")
				*output = o
				continue
			}

			gaucoDtos = append(
				gaucoDtos,
				&spotlikeApp.GetAlbumUseCaseOutputDto{
					ID:          gaucoDto.ID,
					Name:        gaucoDto.Name,
					Artists:     gaucoDto.Artists,
					ReleaseDate: gaucoDto.ReleaseDate,
				},
			)
		}
	}

	clauc := spotlikeApp.NewCheckLikeAlbumUseCase(albumRepo)
	lauc := spotlikeApp.NewLikeAlbumUseCase(albumRepo)
	var likeExecutedAlbums []*spotlikeApp.GetAlbumUseCaseOutputDto
	for _, gaucoDto := range gaucoDtos {
		alreadyLiked, err := clauc.Run(cmd.Context(), gaucoDto.ID)
		if err != nil {
			return err
		}
		if alreadyLiked {
			if err := presenter.Print(os.Stdout, formatter.Blue("⏩ Album "+gaucoDto.Name+" ("+gaucoDto.ID+")"+" released by "+gaucoDto.Artists+" is already liked. skipping...")); err != nil {
				return err
			}
			continue
		}

		if !likeAlbumOps.NoConfirm {
			if answer, err := presenter.RunPrompt(
				"Proceed with liking " + gaucoDto.Name + " (" + gaucoDto.ID + ") ? [y/N]",
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
				o := formatter.Yellow("🚫 Cancelled liking album " + gaucoDto.Name + " (" + gaucoDto.ID + ") ...")
				*output = o
				continue
			}
		}

		if err := lauc.Run(cmd.Context(), gaucoDto.ID); err != nil {
			return err
		}

		likeExecutedAlbums = append(likeExecutedAlbums, gaucoDto)
	}

	if len(likeExecutedAlbums) != 0 {
		f, err := formatter.NewFormatter(likeAlbumOps.Format)
		if err != nil {
			o := formatter.Red("❌ Failed to create a formatter...")
			*output = o
			return err
		}
		o, err := f.Format(likeExecutedAlbums)
		if err != nil {
			return err
		}
		*output = "\n" + o
		if err := presenter.Print(os.Stdout, formatter.Green("✅🤍💿 Successfully liked albums below!")); err != nil {
			return err
		}
	}

	return nil
}

const (
	// likeAlbumHelpTemplate is a help template for the like album command.
	likeAlbumHelpTemplate = `🤍💿 Like albums on Spotify by ID.

You can like tracks on Spotify by ID.

Before using this command,
you need to get the ID of the track you want to like by using the search command.

Also, you can like all albums released by the artist with specifying the ID of the artist with artist flag.
If you specify artist flag, the arguments would be ignored.

` + likeAlbumUsageTemplate
	// likeAlbumUsageTemplate is the usage template of the like album command.
	likeAlbumUsageTemplate = `Usage:
  spotlike like album [flags] [arguments]
  spotlike like al    [flags] [arguments]
  spotlike like a     [flags] [arguments]

Flags:
  -A, --artist  🆔 an ID of the artist to like all albums released by the artist
  --no-confirm  🚫 do not confirm before liking the album
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for album

Arguments:
  ID  🆔 ID of the albums (e.g: "1dGzXXa8MeTCdi0oBbvB1J 6Xy481vVb9vPK4qbCuT9u1")
`
)
