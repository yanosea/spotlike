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
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

const (
	like_artist_help_template = `🤍🎤  Like artist(s) in Spotify by ID.

You must set the args or the flag "id" of artist(s) for like.
If you set both args and flag "id", both will be liked.

Usage:
  spotlike like artist [flags]

Flags:
  -f, --force       like artist(s) without confirming
  -h, --help        help for artist
  -i, --id string   ID of the artist(s) for like
  -v, --verbose     print verbose output
`
	like_artist_use   = "artist"
	like_artist_short = "🤍🎤 Like artist(s) in Spotify by ID."
	like_artist_long  = `🤍🎤 Like artist(s) in Spotify by ID.

You must set the args or the flag "id" of artist(s) for like.
If you set both args and flag "id", both will be liked.`
	like_artist_flag_id                                    = "id"
	like_artist_shorthand_id                               = "i"
	like_artist_flag_description_id                        = "ID of the artist(s) for like"
	like_artist_flag_force                                 = "force"
	like_artist_shorthand_force                            = "f"
	like_artist_flag_description_force                     = "like artist(s) without confirming"
	like_artist_flag_verbose                               = "verbose"
	like_artist_shorthand_verbose                          = "v"
	like_artist_flag_description_verbose                   = "print verbose output"
	like_artist_error_message_template_id_not_artist       = "The ID you passed [%s] is not artist..."
	like_artist_message_template_like_artist_already_liked = "✨ %s already liked!\t:\t[%s]"
	like_artist_message_template_like_artist_refused       = "❌ Like %s refused!\t:\t[%s]"
	like_artist_message_template_like_artist_succeeded     = "✅ Like %s succeeded!\t:\t[%s]"
)

type likeArtistOption struct {
	Client *spotify.Client

	Args    []string
	Id      string
	Force   bool
	Verbose bool

	Out    io.Writer
	ErrOut io.Writer
}

func newLikeArtistCommand(globalOption *GlobalOption) *cobra.Command {
	o := &likeArtistOption{}

	cmd := &cobra.Command{
		Use:   like_artist_use,
		Short: like_artist_short,
		Long:  like_artist_long,
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

			return o.likeArtist()
		},
	}

	cmd.PersistentFlags().StringVarP(&o.Id, like_artist_flag_id, like_artist_shorthand_id, "", like_artist_flag_description_id)
	cmd.PersistentFlags().BoolVarP(&o.Force, like_artist_flag_force, like_artist_shorthand_force, false, like_artist_flag_description_force)
	cmd.PersistentFlags().BoolVarP(&o.Verbose, like_artist_flag_verbose, like_artist_shorthand_verbose, false, like_artist_flag_description_verbose)

	o.Out = globalOption.Out
	o.ErrOut = globalOption.ErrOut
	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	cmd.SetHelpTemplate(like_artist_help_template)

	return cmd
}

func (o *likeArtistOption) likeArtist() error {
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
		// first, search the artist by the id
		searchResult, err := api.SearchById(o.Client, id)
		if err != nil {
			return err
		}
		if searchResult.Type != spotify.SearchTypeArtist {
			// if the search result was not artist
			return errors.New(fmt.Sprintf(like_artist_error_message_template_id_not_artist, id))
		}
		likeResults = api.LikeArtistById(o.Client, searchResult.Id, o.Force)
	}
	// print the like result
	o.printLikeArtistResult(likeResults)

	return nil
}

func (o *likeArtistOption) printLikeArtistResult(likeResults []*api.LikeResult) {
	for _, result := range likeResults {
		if (result.Refused || result.AlreadyLiked) && len(likeResults) != 1 && !o.Verbose {
			// if the result was refused or already liked
			// and the result was not one,
			// and the verbose option was not set
			// skip the result
			continue
		}
		if result.Error != nil {
			util.PrintlnWithWriter(o.ErrOut, FormatLikeResultError(result.Error))
		}
		util.PrintlnWithWriter(o.Out, formatLikeArtistResult(result))
	}
}

func formatLikeArtistResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		return color.BlueString(fmt.Sprintf(like_artist_message_template_like_artist_already_liked, util.STRING_ARTIST, result.ArtistNames))
	} else if result.Refused {
		return color.YellowString(fmt.Sprintf(like_artist_message_template_like_artist_refused, util.STRING_ARTIST, result.ArtistNames))
	} else {
		return color.GreenString(fmt.Sprintf(like_artist_message_template_like_artist_succeeded, util.STRING_ARTIST, result.ArtistNames))
	}
}
