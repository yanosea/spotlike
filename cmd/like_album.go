package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/library/api"
	"github.com/yanosea/spotlike/app/library/auth"
	"github.com/yanosea/spotlike/app/library/utility"
	"github.com/yanosea/spotlike/app/proxy/cobra"
	"github.com/yanosea/spotlike/app/proxy/color"
	"github.com/yanosea/spotlike/app/proxy/io"
	"github.com/yanosea/spotlike/cmd/constant"
)

type likeAlbumOption struct {
	Out     ioproxy.WriterInstanceInterface
	ErrOut  ioproxy.WriterInstanceInterface
	Args    []string
	Utility utility.UtilityInterface
	Id      string
	Force   bool
	Verbose bool
}

func NewLikeAlbumCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &likeAlbumOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    g.Args,
		Utility: g.Utility,
	}

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.LIKE_ALBUM_USE
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.likeAlbumRunE

	cmd.PersistentFlags().StringVarP(
		&o.Id,
		constant.LIKE_ALBUM_FLAG_ID,
		constant.LIKE_ALBUM_FLAG_ID_SHORTHAND,
		"",
		constant.LIKE_ALBUM_FLAG_ID_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.Force,
		constant.LIKE_ALBUM_FLAG_FORCE,
		constant.LIKE_ALBUM_FLAG_FORCE_SHORTHAND,
		false,
		constant.LIKE_ALBUM_FLAG_FORCE_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.Verbose,
		constant.LIKE_ALBUM_FLAG_VERBOSE,
		constant.LIKE_ALBUM_FLAG_VERBOSE_SHORTHAND,
		false,
		constant.LIKE_ALBUM_FLAG_VERBOSE_DESCRIPTION,
	)

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.LIKE_ALBUM_HELP_TEMPLATE)

	return cmd
}

func (o *likeAlbumOption) likeAlbumRunE(_ *cobra.Command, _ []string) error {
	if o.Client == nil {
		if !auth.IsEnvsSet() {
			auth.SetAuthInfo()
		}
		client, err := auth.Authenticate(g.Out)
		if err != nil {
			return err
		}
		o.Client = client
	}

	return o.likeAlbum()
}

func (o *likeAlbumOption) likeAlbum() error {
	// set the the id(s) from the args and the flag
	var ids []string
	for _, arg := range o.Args {
		for _, a := range strings.Fields(arg) {
			ids = append(ids, a)
		}
	}
	for _, a := range strings.Fields(o.Id) {
		ids = append(ids, a)
	}
	if len(ids) == 0 {
		// if the id(s) was not set
		return errors.New(like_error_message_args_or_flag_id_required)
	}
	var likeResults []*api.LikeResult
	for _, id := range ids {
		// first, search the content by the id
		searchResult, err := api.SearchById(o.Client, id)
		if err != nil {
			return err
		}
		switch searchResult.Type {
		case spotify.SearchTypeArtist:
			// confirm like all albums by the artist
			answer := "y"
			if !o.Force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf(like_album_input_label_template_all_album_by_artist, searchResult.ArtistNames),
					IsConfirm: true,
				}

				input, err := prompt.Run()
				if err != nil || input == "" {
					answer = "n"
				}
			}
			if answer == "y" {
				// like all albums released by the artist
				likeResults = api.LikeAllAlbumsReleasedByArtistById(o.Client, searchResult.Id, o.Force)
			}
		case spotify.SearchTypeAlbum:
			// like the album
			likeResults = api.LikeAlbumById(o.Client, searchResult.Id, o.Force)
		default:
			// if the id was not artist or album
			return errors.New(fmt.Sprintf(constant.LIKE_ALBUM_ERROR_MESSAGE_TEMPLATE_ID_NOT_ARTIST_ALBUM, id))
		}
		// print the like result
		o.printLikeAlbumResult(likeResults)
	}

	return nil
}

func (o *likeAlbumOption) printLikeAlbumResult(likeResults []*api.LikeResult) {
	for _, result := range likeResults {
		if (result.Refused || result.AlreadyLiked) && len(likeResults) != 1 && !o.Verbose {
			// if the result was refused or already liked
			// and the result was not one,
			// and the verbose option was not set
			// skip the result
			continue
		}
		if result.Error != nil {
			o.Utility.PrintlnWithWriter(o.ErrOut, formatLikeResultError(result.Error))
		}
		o.Utility.PrintlnWithWriter(o.Out, formatLikeAlbumResult(result))
	}
}

func formatLikeAlbumResult(result *api.LikeResult) string {
	colorProxy := colorproxy.New()
	if result.AlreadyLiked {
		return colorProxy.BlueString(fmt.Sprintf(like_album_message_template_like_album_already_liked, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	} else if result.Refused {
		return colorProxy.YellowString(fmt.Sprintf(like_album_message_template_like_album_refused, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	} else {
		return colorProxy.GreenString(fmt.Sprintf(like_album_message_template_like_album_succeeded, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	}
}
