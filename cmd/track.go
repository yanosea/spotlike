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
	like_track_use   = "track"
	like_track_short = "Like track(s) in Spotify by ID."
	like_track_long  = `Like track(s) in Spotify by ID.

You must set the args or the flag "id" of track(s) for like.
If you set both args and flag "id", both tracks will be liked.

You can also set the flag "album" or "artist".
If you set the flag "artist" and pass the artist ID, spotlike will like all tracks released by the artist.
If you set the flag "album" and pass the album ID, spotlike will like all tracks included in the album.`
	like_track_flag_id                                = "id"
	like_track_shorthand_id                           = "i"
	like_track_flag_description_id                    = "ID of the track(s) or the artist(s) or the album for like"
	like_track_flag_force                             = "force"
	like_track_shorthand_force                        = "f"
	like_track_flag_description_force                 = "like track(s) without confirming"
	like_track_flag_verbose                           = "verbose"
	like_track_shorthand_verbose                      = "v"
	like_track_flag_description_verbose               = "print verbose output"
	like_track_flag_artist                            = "artist"
	like_track_flag_description_artist                = "like all albums released by the artist"
	like_track_flag_album                             = "album"
	like_track_flag_description_album                 = "like all tracks included in the album"
	like_track_error_message_args_or_flag_id_required = `The arguments or the flag "id" is required...`
	like_track_error_message_template_id_not_track    = `The ID you passed [%s] is not the track ID.
You have to pass the track ID for like the track(s) or set flag "artist" or "album" and pass the artist ID or the album ID.`
	like_track_error_message_template_flag_artist_invalid = `You set the flag "artist" but the ID you passed [%s] is not the artist ID.
You have to pass the artist ID for like all tracks released by the artist.`
	like_track_error_message_template_flag_album_invalid = `You set the flag "album" but the ID you passed [%s] is not the album ID.
You have to pass the artist ID for like all included in the track.`
	like_track_confirm_message_template_all_track_by_artist = "Do you execute like all tracks by [%s]"
	like_track_confirm_message_template_all_track_in_album  = "Do you execute like all tracks in [%s]"
	like_track_message_template_like_track_already_liked    = "%s in [%s] by [%s] already liked!\t:\t[%s]"
	like_track_message_template_like_track_refused          = "Like %s in [%s] by [%s] refused!\t:\t[%s]"
	like_track_message_template_like_track_succeeded        = "Like %s in [%s] by [%s] succeeded!\t:\t[%s]"
)

type likeTrackOption struct {
	Client *spotify.Client

	Args    []string
	Id      string
	Level   string
	Force   bool
	Verbose bool
	Artist  bool
	Album   bool

	Out    io.Writer
	ErrOut io.Writer
}

func newLikeTrackCommand(globalOption *GlobalOption) *cobra.Command {
	o := &likeTrackOption{}

	cmd := &cobra.Command{
		Use:   like_track_use,
		Short: like_track_short,
		Long:  like_track_long,
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

			return o.likeTrack()
		},
	}

	cmd.PersistentFlags().StringVarP(&o.Id, like_track_flag_id, like_track_shorthand_id, "", like_track_flag_description_id)
	cmd.PersistentFlags().BoolVarP(&o.Force, like_track_flag_force, like_track_shorthand_force, false, like_track_flag_description_force)
	cmd.PersistentFlags().BoolVarP(&o.Verbose, like_track_flag_verbose, like_track_shorthand_verbose, false, like_track_flag_description_verbose)
	cmd.PersistentFlags().BoolVarP(&o.Artist, like_track_flag_artist, "", false, like_track_flag_description_artist)
	cmd.PersistentFlags().BoolVarP(&o.Artist, like_track_flag_album, "", false, like_track_flag_description_album)
	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
}

func (o *likeTrackOption) likeTrack() error {
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
		if searchResult.Type != spotify.SearchTypeTrack {
			if o.Artist && searchResult.Type != spotify.SearchTypeArtist {
				// if the search result was not track
				// or the search result was not artist and the artist flag was not set
				return errors.New(fmt.Sprintf(like_track_error_message_template_flag_artist_invalid, id))
			}
			if o.Album && searchResult.Type != spotify.SearchTypeAlbum {
				// if the search result was not track
				// or the search result was not album and the album flag was not set
				return errors.New(fmt.Sprintf(like_track_error_message_template_flag_album_invalid, id))
			}
		}
		switch searchResult.Type {
		case spotify.SearchTypeArtist:
			// confirm like all tracks by the artist
			answer := "y"
			if !o.Force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf(like_track_confirm_message_template_all_track_by_artist, searchResult.ArtistNames),
					IsConfirm: true,
				}

				input, err := prompt.Run()
				if err != nil || input == "" {
					answer = "n"
				}
			}
			if answer == "y" {
				likeResults = api.LikeAllTracksReleasedByArtistById(o.Client, searchResult.Id, o.Force)
			}
		case spotify.SearchTypeAlbum:
			// confirm like all tracks included in the album
			answer := "y"
			if !o.Force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf(like_track_confirm_message_template_all_track_in_album, searchResult.AlbumName),
					IsConfirm: true,
				}

				input, err := prompt.Run()
				if err != nil || input == "" {
					answer = "n"
				}
			}
			if answer == "y" {
				likeResults = api.LikeAllTracksInAlbumById(o.Client, searchResult.Id, o.Force)
			}
		case spotify.SearchTypeTrack:
			likeResults = api.LikeTrackById(o.Client, searchResult.Id, o.Force)
		}
		o.printLikeTrackResult(likeResults)
	}

	return nil
}

func (o *likeTrackOption) printLikeTrackResult(likeResults []*api.LikeResult) {
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
		util.PrintlnWithWriter(o.Out, formatLikeTrackResult(result))
	}
}

func formatLikeTrackResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		return color.BlueString(fmt.Sprintf(like_track_message_template_like_track_already_liked, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	} else if result.Refused {
		return color.YellowString(fmt.Sprintf(like_track_message_template_like_track_refused, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	} else {
		return color.GreenString(fmt.Sprintf(like_track_message_template_like_track_succeeded, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	}
}

func formatLikeResultError(error error) string {
	return color.RedString(fmt.Sprintf("Error:\n  %s", error))
}
