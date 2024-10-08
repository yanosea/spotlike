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

type likeTrackOption struct {
	Out     ioproxy.WriterInstanceInterface
	ErrOut  ioproxy.WriterInstanceInterface
	Args    []string
	Utility utility.UtilityInterface
	Id      string
	Force   bool
	Level   string
	Verbose bool
}

func NewLikeTrackCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &likeTrackOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    g.Args,
		Utility: g.Utility,
	}

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.LIKE_TRACK_USE
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.likeTrackRunE

	cmd.PersistentFlags().StringVarP(&o.Id,
		constant.LIKE_TRACK_FLAG_ID,
		constant.LIKE_TRACK_FLAG_ID_SHORTHAND,
		"",
		constant.LIKE_TRACK_FLAG_ID_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Force,
		constant.LIKE_TRACK_FLAG_FORCE,
		constant.LIKE_TRACK_FLAG_FORCE_SHORTHAND,
		false,
		constant.LIKE_TRACK_FLAG_FORCE_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Verbose,
		constant.LIKE_TRACK_FLAG_VERBOSE,
		constant.LIKE_TRACK_FLAG_VERBOSE_SHORTHAND,
		false,
		constant.LIKE_TRACK_FLAG_VERBOSE_DESCRIPTION,
	)

	o.Out = g.Out
	o.ErrOut = g.ErrOut
	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	cmd.SetHelpTemplate(constant.LIKE_TRACK_HELP_TEMPLATE)

	return cmd
}

func (o *likeTrackOption) likeTrackRunE(_ *cobra.Command, _ []string) error {
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

	return o.likeTrack()
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
		switch searchResult.Type {
		case spotify.SearchTypeArtist:
			// confirm like all tracks by the artist
			answer := "y"
			if !o.Force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf(like_track_input_label_template_all_track_by_artist, searchResult.ArtistNames),
					IsConfirm: true,
				}

				input, err := prompt.Run()
				if err != nil || input == "" {
					answer = "n"
				}
			}
			if answer == "y" {
				// like all tracks released by the artist
				likeResults = api.LikeAllTracksReleasedByArtistById(o.Client, searchResult.Id, o.Force)
			}
		case spotify.SearchTypeAlbum:
			// confirm like all tracks included in the album
			answer := "y"
			if !o.Force {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf(like_track_input_label_template_all_track_in_album, searchResult.AlbumName),
					IsConfirm: true,
				}

				input, err := prompt.Run()
				if err != nil || input == "" {
					answer = "n"
				}
			}
			if answer == "y" {
				// like all tracks included in the album
				likeResults = api.LikeAllTracksInAlbumById(o.Client, searchResult.Id, o.Force)
			}
		case spotify.SearchTypeTrack:
			// like the track
			likeResults = api.LikeTrackById(o.Client, searchResult.Id, o.Force)
		default:
			// if the id was not artist, album or track
			return errors.New(fmt.Sprintf(like_track_error_message_template_id_not_artist_album_track, id))
		}
		// print the like result
		o.printLikeTrackResult(likeResults)
	}

	return nil
}

func (o *likeTrackOption) printLikeTrackResult(likeResults []*api.LikeResult) {
	for _, result := range likeResults {
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
	colorProxy := colorproxy.New()
	if result.AlreadyLiked {
		return colorProxy.BlueString(fmt.Sprintf(like_track_message_template_like_track_already_liked, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	} else if result.Refused {
		return colorProxy.YellowString(fmt.Sprintf(like_track_message_template_like_track_refused, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	} else {
		return colorProxy.GreenString(fmt.Sprintf(like_track_message_template_like_track_succeeded, util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName))
	}
}

func formatLikeResultError(error error) string {
	colorProxy := colorproxy.New()
	return colorProxy.RedString(fmt.Sprintf("Error:\n  %s", error))
}
