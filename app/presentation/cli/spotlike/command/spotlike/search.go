package spotlike

import (
	c "github.com/spf13/cobra"

	spotlikeApp "github.com/yanosea/spotlike/app/application/spotlike"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/api"
	"github.com/yanosea/spotlike/app/infrastructure/spotify/repository"
	"github.com/yanosea/spotlike/app/presentation/cli/spotlike/formatter"

	"github.com/yanosea/spotlike/pkg/proxy"
)

// SearchOptions provides the options for the search command.
type SearchOptions struct {
	Artist bool
	Album  bool
	Track  bool
	Max    int
	Format string
}

var (
	// searchOps is a variable to store the search options with the default values for injecting the dependencies in testing.
	searchOps = SearchOptions{
		Artist: false,
		Album:  false,
		Track:  false,
		Max:    10,
		Format: "table",
	}
)

// NewSearchCommand returns a new instance of the search command.
func NewSearchCommand(
	cobra proxy.Cobra,
	authCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("search")
	cmd.SetAliases([]string{"se", "s"})
	cmd.SetUsageTemplate(searchUsageTemplate)
	cmd.SetHelpTemplate(searchHelpTemplate)
	cmd.SetSilenceErrors(true)
	cmd.Flags().BoolVarP(
		&searchOps.Artist,
		"artist",
		"A",
		false,
		"🎤 search for artists",
	)
	cmd.Flags().BoolVarP(
		&searchOps.Album,
		"album",
		"a",
		false,
		"💿 search for albums",
	)
	cmd.Flags().BoolVarP(
		&searchOps.Track,
		"track",
		"t",
		false,
		"🎵 search for tracks",
	)
	cmd.Flags().IntVarP(
		&searchOps.Max,
		"max",
		"m",
		10,
		"🔢 maximum number of search results",
	)
	cmd.Flags().StringVarP(
		&searchOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g: \"plain\")",
	)
	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runSearch(cmd, authCmd, output, args)
		},
	)

	return cmd
}

// runSearch runs the search command.
func runSearch(cmd *c.Command, authCmd proxy.Command, output *string, args []string) error {
	if len(args) == 0 {
		o := formatter.Yellow("⚡ No query arguments specified...")
		*output = o
		return nil
	}

	if !searchOps.Artist && !searchOps.Album && !searchOps.Track {
		o := formatter.Yellow("⚡ No search type specified...")
		*output = o
		return nil
	}

	if searchOps.Artist && searchOps.Album ||
		searchOps.Artist && searchOps.Track ||
		searchOps.Album && searchOps.Track {
		o := formatter.Yellow("⚡ Multiple search types is not supported...")
		*output = o
		return nil
	}

	if searchOps.Max < 1 {
		o := formatter.Yellow("⚡ Invalid search result number...")
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

	var dtos any
	if searchOps.Artist {
		artistRepo := repository.NewArtistRepository()
		sAuc := spotlikeApp.NewSearchArtistUseCase(artistRepo)
		sAoDtos, err := sAuc.Run(cmd.Context(), args, searchOps.Max)
		if err != nil {
			return err
		}

		if len(sAoDtos) == 0 {
			o := formatter.Yellow("⚡ No artists found...")
			*output = o
			return nil
		}
		dtos = sAoDtos
	}

	if searchOps.Album {
		albumRepo := repository.NewAlbumRepository()
		sauc := spotlikeApp.NewSearchAlbumUseCase(albumRepo)
		saoDtos, err := sauc.Run(cmd.Context(), args, searchOps.Max)
		if err != nil {
			return err
		}

		if len(saoDtos) == 0 {
			o := formatter.Yellow("⚡ No albums found...")
			*output = o
			return nil
		}
		dtos = saoDtos
	}

	if searchOps.Track {
		trackRepo := repository.NewTrackRepository()
		stuc := spotlikeApp.NewSearchTrackUseCase(trackRepo)
		stoDtos, err := stuc.Run(cmd.Context(), args, searchOps.Max)
		if err != nil {
			return err
		}

		if len(stoDtos) == 0 {
			o := formatter.Yellow("⚡ No tracks found...")
			*output = o
			return nil
		}
		dtos = stoDtos
	}

	f, err := formatter.NewFormatter(searchOps.Format)
	if err != nil {
		o := formatter.Red("❌ Failed to create a formatter...")
		*output = o
		return err
	}
	o, err := f.Format(dtos)
	if err != nil {
		return err
	}
	*output = "\n" + o

	return nil
}

const (
	// searchHelpTemplate is the help template of the search command.
	searchHelpTemplate = `🔍 Search for the ID of content in Spotify.

You can search for artists, albums, and tracks by specifying the type of content to search for.

If you want to search artists, specify the "-A" or "--artist" option.
If you want to search albums, specify the "-a" or "--album" option.
If you want to search tracks, specify the "-t" or "--track" option.

Also, you can specify the maximum number of search results by specifying the "-m" or "--max" option.

` + searchUsageTemplate
	// searchUsageTemplate is the usage template of the search command.
	searchUsageTemplate = `Usage:
  spotlike search [flags] [arguments]
  spotlike se     [flags] [arguments]
  spotlike s      [flags] [arguments]

Flags:
  -A, --artist  🎤 search for artists
  -a, --album   💿 search for albums
  -t, --track   🎵 search for tracks
  -m, --max     🔢 maximum number of search results (default 10)
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for search

Arguments:
  keywords  🔡 search content by keywords (multiple keywords are separated by a space)
`
)
