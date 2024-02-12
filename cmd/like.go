package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/yanosea/spotlike/api"
	"github.com/yanosea/spotlike/auth"
	"github.com/yanosea/spotlike/util"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

const (
	like_use   = "like"
	like_short = "Like content in Spotify by ID."
	like_long  = `Like content in Spotify by ID.

You must set the args or the flag "id" to content(s) for like.
If you set both args and flag "id", both contents will be liked.

You can like content(s) using the flag "level" below :
  * artist (If you pass the artist ID, spotlike will like the artist.)
  * album  (If you pass the artist ID, spotlike will like all albums released by the artist.)
           (If you pass the album ID, spotlike will like the album.)
  * track  (If you pass the artist ID, spotlike will like all tracks released by the artist.)
           (If you pass the album ID, spotlike will like all tracks contained in the album.)
           (If you pass the track ID, spotlike will like the track.)`
	like_flag_id                                      = "id"
	like_shorthand_id                                 = "i"
	like_flag_description_id                          = "ID of the content for like"
	like_flag_level                                   = "level"
	like_shorthand_level                              = "l"
	like_flag_description_level                       = "level for like"
	like_flag_force                                   = "force"
	like_shorthand_force                              = "f"
	like_flag_description_force                       = "like without confirming"
	like_error_message_args_or_flag_id_required       = `The arguments or the flag "id" is required...`
	like_error_message_flag_level_invalid_artist      = "You passed the artist ID, so you have to pass 'artist', 'album', or 'track' for the level option. Or you should not specify the level option to like the artist."
	like_error_message_flag_level_invalid_album       = "You passed the album ID, so you have to pass 'album' or 'track' for the level option. Or you should not specify the level option to like the album."
	like_error_message_flag_level_invalid_track       = "You passed the track ID, so you have to pass 'track' for the level option. Or you should not specify the level option to like the track."
	like_error_message_something_wrong_with_searching = "Search result is wrong..."
	message_like_artist_already_liked                 = "%s already liked!\t:\t[%s]"
	message_like_artist_refused                       = "Like %s refused!\t:\t[%s]"
	message_like_artist_succeeded                     = "Like %s succeeded!\t:\t[%s]"
	message_like_album_already_liked                  = "%s by [%s] already liked!\t:\t[%s]"
	message_like_album_refused                        = "Like %s by [%s] refused!\t:\t[%s]"
	message_like_album_succeeded                      = "Like %s by [%s] succeeded!\t:\t[%s]"
	message_like_track_already_liked                  = "%s in [%s] by [%s] already liked!\t:\t[%s]"
	message_like_track_refused                        = "Like %s in [%s] by [%s] refused!\t:\t[%s]"
	message_like_track_succeeded                      = "Like %s in [%s] by []%s] succeeded!\t:\t[%s]"
)

type likeOption struct {
	NoColor bool
	Args    []string
	Id      string
	Level   string
	Force   bool
	Client  *spotify.Client
	Out     io.Writer
	ErrOut  io.Writer
}

func newLikeCommand(globalOption *GlobalOption) *cobra.Command {
	o := &likeOption{}
	o.Out = globalOption.Out
	o.ErrOut = globalOption.ErrOut

	cmd := &cobra.Command{
		Use:   like_use,
		Short: like_short,
		Long:  like_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Client = globalOption.Client
			if o.Client == nil {
				client, err := auth.GetClient()
				if err != nil {
					return err
				}
				o.Client = client
			}

			cmd.Flags().StringVarP(&o.Id, like_flag_id, like_shorthand_id, "", like_flag_description_id)
			cmd.Flags().StringVarP(&o.Level, like_flag_level, like_shorthand_level, "", like_flag_description_level)
			cmd.Flags().BoolVarP(&o.Force, like_flag_force, like_shorthand_force, false, like_flag_description_force)

			o.Args = args

			return o.like()
		},
	}

	return cmd
}

func (o *likeOption) like() error {
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
		return errors.New(like_error_message_args_or_flag_id_required)
	}

	for _, id := range ids {
		searchResult, err := api.SearchById(o.Client, id)
		if err != nil {
			return err
		}
		var likeResults []*api.LikeResult
		switch searchResult.Type {
		case spotify.SearchTypeArtist:
			if strings.ToLower(o.Level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeArtist] || o.Level == "" {
				likeResults = api.LikeArtistById(o.Client, searchResult.Id, o.Force)
			} else if strings.ToLower(o.Level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeAlbum] {
				likeResults = api.LikeAllAlbumsReleasedByArtistById(o.Client, searchResult.Id, o.Force)
			} else if strings.ToLower(o.Level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeTrack] {
				likeResults = api.LikeAllTracksReleasedByArtistById(o.Client, searchResult.Id, o.Force)
			} else {
				return errors.New(like_error_message_flag_level_invalid_artist)
			}
		case spotify.SearchTypeAlbum:
			if strings.ToLower(o.Level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeAlbum] || o.Level == "" {
				likeResults = api.LikeAlbumById(o.Client, searchResult.Id, o.Force)
			} else if strings.ToLower(o.Level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeTrack] {
				likeResults = api.LikeAllTracksInAlbumById(o.Client, searchResult.Id, o.Force)
			} else {
				return errors.New(like_error_message_flag_level_invalid_album)
			}
		case spotify.SearchTypeTrack:
			if strings.ToLower(o.Level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeTrack] || o.Level == "" {
				likeResults = api.LikeTrackById(o.Client, searchResult.Id, o.Force)
			} else {
				return errors.New(like_error_message_flag_level_invalid_track)
			}
		default:
			return errors.New(like_error_message_something_wrong_with_searching)
		}
		o.printLikeResult(likeResults)
	}

	return nil
}

func (o *likeOption) printLikeResult(likeResults []*api.LikeResult) {
	for _, result := range likeResults {
		if result.Error != nil {
			util.PrintlnWithWriter(o.ErrOut, formatLikeResultError(result.Error))
		}
		if result.Type == spotify.SearchTypeArtist {
			util.PrintlnWithWriter(o.Out, formatLikeArtistResult(result))
		}
		if result.Type == spotify.SearchTypeAlbum {
			util.PrintlnWithWriter(o.Out, formatLikeAlbumResult(result))
		}
		if result.Type == spotify.SearchTypeTrack {
			util.PrintlnWithWriter(o.Out, formatLikeTrackResult(result))
		}
	}
}

func formatLikeArtistResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		return color.YellowString(fmt.Sprintf(message_like_artist_already_liked, util.STRING_ARTIST, result.ArtistNames))
	} else if result.Refused {
		return color.YellowString(fmt.Sprintf(message_like_artist_refused, util.STRING_ARTIST, result.ArtistNames))
	} else {
		return color.GreenString(fmt.Sprintf(message_like_artist_succeeded, util.STRING_ARTIST, result.ArtistNames))
	}
}

func formatLikeAlbumResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		return color.YellowString(fmt.Sprintf(message_like_album_already_liked, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	} else if result.Refused {
		return color.YellowString(fmt.Sprintf(message_like_album_refused, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	} else {
		return color.GreenString(fmt.Sprintf(message_like_album_succeeded, util.STRING_ALBUM, result.ArtistNames, result.AlbumName))
	}
}

func formatLikeTrackResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		return color.YellowString(fmt.Sprintf(message_like_track_already_liked, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	} else if result.Refused {
		return color.YellowString(fmt.Sprintf(message_like_track_refused, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	} else {
		return color.GreenString(fmt.Sprintf(message_like_track_succeeded, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	}
}

func formatLikeResultError(error error) string {
	return color.RedString(fmt.Sprintf("Error:\n  %s", error))
}
