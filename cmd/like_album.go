package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/yanosea/spotlike/api"
	"github.com/yanosea/spotlike/auth"
	"github.com/yanosea/spotlike/help"
	"github.com/yanosea/spotlike/util"

	// https://github.com/fatih/color
	"github.com/fatih/color"
	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

const (
	like_album_use   = "album"
	like_album_short = "🤍💿 Like album(s) in Spotify by ID."
	like_album_long  = `🤍💿 Like album(s) in Spotify by ID.

You must set the args or the flag "id" of album(s) or artist(s) for like.
If you set both args and flag "id", both will be liked.

If you pass the artist ID, spotlike will like all albums released by the artist.`
	like_album_flag_id                                    = "id"
	like_album_shorthand_id                               = "i"
	like_album_flag_description_id                        = "ID of the album(s) or the artist(s) for like"
	like_album_flag_force                                 = "force"
	like_album_shorthand_force                            = "f"
	like_album_flag_description_force                     = "like album(s) without confirming"
	like_album_flag_verbose                               = "verbose"
	like_album_shorthand_verbose                          = "v"
	like_album_flag_description_verbose                   = "print verbose output"
	like_album_error_message_template_id_not_artist_album = "The ID you passed [%s] is not album or artist."
	like_album_input_label_template_all_album_by_artist   = "Do you execute like all albums by [%s]"
	like_album_message_template_like_album_already_liked  = "%s by [%s] already liked!\t:\t[%s]"
	like_album_message_template_like_album_refused        = "Like %s by [%s] refused!\t:\t[%s]"
	like_album_message_template_like_album_succeeded      = "Like %s by [%s] succeeded!\t:\t[%s]"
)

type likeAlbumOption struct {
	Client *spotify.Client

	Args    []string
	Id      string
	Force   bool
	Verbose bool

	Out    io.Writer
	ErrOut io.Writer
}

func newLikeAlbumCommand(globalOption *GlobalOption) *cobra.Command {
	o := &likeAlbumOption{}

	cmd := &cobra.Command{
		Use:   like_album_use,
		Short: like_album_short,
		Long:  like_album_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Client = globalOption.Client
			if o.Client == nil {
				if !auth.IsEnvsSet() {
					auth.SetAuthInfo()
				}
				client, err := auth.Authenticate(globalOption.Out)
				if err != nil {
					return err
				}
				o.Client = client
			}
			o.Args = args
			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut

			return o.likeAlbum()
		},
	}

	cmd.PersistentFlags().StringVarP(&o.Id, like_album_flag_id, like_album_shorthand_id, "", like_album_flag_description_id)
	cmd.PersistentFlags().BoolVarP(&o.Force, like_album_flag_force, like_album_shorthand_force, false, like_album_flag_description_force)
	cmd.PersistentFlags().BoolVarP(&o.Verbose, like_album_flag_verbose, like_album_shorthand_verbose, false, like_album_flag_description_verbose)

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	cmd.SetHelpTemplate(help.LIKE_ALBUM_HELP_TEMPLATE)

	return cmd
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
			return errors.New(fmt.Sprintf(like_album_error_message_template_id_not_artist_album, id))
		}
		// print the like result
		o.printLikeAlbumResult(likeResults)
	}

	return nil
}

func (o *likeAlbumOption) printLikeAlbumResult(likeResults []*api.LikeResult) {
	for index, result := range likeResults {
		if index == 0 {
			// if the result was first, add blank line
			util.PrintlnWithWriter(o.Out, "")
		}
		if (result.Refused || result.AlreadyLiked) && len(likeResults) != 1 && !o.Verbose {
			// if the result was refused or already liked
			// and the result was not one,
			// and the verbose option was not set
			// skip the result
			continue
		}
		if result.Error != nil {
			util.PrintlnWithWriter(o.ErrOut, formatLikeResultError(result.Error))
		}
		util.PrintlnWithWriter(o.Out, formatLikeAlbumResult(result))
	}
}

func formatLikeAlbumResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		return color.BlueString(fmt.Sprintf(like_album_message_template_like_album_already_liked, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	} else if result.Refused {
		return color.YellowString(fmt.Sprintf(like_album_message_template_like_album_refused, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	} else {
		return color.GreenString(fmt.Sprintf(like_album_message_template_like_album_succeeded, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	}
}
