package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/yanosea/spotlike/api"
	"github.com/yanosea/spotlike/auth"
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
	like_album_short = "Like album(s) in Spotify by ID."
	like_album_long  = `Like album(s) in Spotify by ID.

You must set the args or the flag "id" of album(s) for like.
If you set both args and flag "id", both albums will be liked.

You can also set the flag "artist".
If you set the flag "artist" and pass the artist ID, spotlike will like all albums released by the artist.`
	like_album_flag_id                                           = "id"
	like_album_shorthand_id                                      = "i"
	like_album_flag_description_id                               = "ID of the album(s) or the artist(s) for like"
	like_album_flag_force                                        = "force"
	like_album_shorthand_force                                   = "f"
	like_album_flag_description_force                            = "like album(s) without confirming"
	like_album_flag_verbose                                      = "verbose"
	like_album_shorthand_verbose                                 = "v"
	like_album_flag_artist                                       = "artist"
	like_album_flag_description_artist                           = "like all albums released by the artist"
	like_album_flag_description_verbose                          = "print verbose output"
	like_album_error_message_template_flag_artist_invalid_artist = `The ID you passed [%s] is artist ID but you did not set the flag "artist".
You have to set the flag "artist" ID with setting the flag "artist".`
	like_album_error_message_template_flag_artist_invalid_album = `You set the flag "artist" but the ID you passed [%s] is album ID.
You have to pass the album ID without setting the flag "artist".`
	like_album_error_message_template_id_not_album_artist = `The ID you passed [%s] is neither the album ID nor the artist ID.
You have to pass the album ID for like the album(s).
You have to pass the artist Id for like the all albums released by the artist with setting the flag "artist".`
	like_album_confirm_message_template_all_album_by_artist = "Do you execute like all albums by [%s]"
	like_album_message_template_like_album_already_liked    = "%s by [%s] already liked!\t:\t[%s]"
	like_album_message_template_like_album_refused          = "Like %s by [%s] refused!\t:\t[%s]"
	like_album_message_template_like_album_succeeded        = "Like %s by [%s] succeeded!\t:\t[%s]"
)

type likeAlbumOption struct {
	Client *spotify.Client

	Args    []string
	Id      string
	Force   bool
	Verbose bool
	Artist  bool

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
	cmd.PersistentFlags().BoolVarP(&o.Artist, like_album_flag_artist, "", false, like_album_flag_description_artist)
	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

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
			if !o.Artist {
				// if the search result was artist and the flag "artist" was not set
				return errors.New(fmt.Sprintf(like_album_error_message_template_flag_artist_invalid_artist, id))
			}
		case spotify.SearchTypeAlbum:
			// if the search result was album and the flag "artist" was set
			if o.Artist {
				return errors.New(fmt.Sprintf(like_album_error_message_template_flag_artist_invalid_album, id))
			}
		default:
			// if the search result was not album and not artist
			return errors.New(fmt.Sprintf(like_album_error_message_template_id_not_album_artist, id))
		}
		switch searchResult.Type {
		case spotify.SearchTypeArtist:
			// confirm like all albums by the artist
			answer := "y"
			if !o.Force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf(like_album_confirm_message_template_all_album_by_artist, searchResult.ArtistNames),
					IsConfirm: true,
				}

				input, err := prompt.Run()
				if err != nil || input == "" {
					answer = "n"
				}
			}
			if answer == "y" {
				likeResults = api.LikeAllAlbumsReleasedByArtistById(o.Client, searchResult.Id, o.Force)
			}
		case spotify.SearchTypeAlbum:
			likeResults = api.LikeAlbumById(o.Client, searchResult.Id, o.Force)
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
