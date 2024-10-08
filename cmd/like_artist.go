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

type likeArtistOption struct {
	Out     ioproxy.WriterInstanceInterface
	ErrOut  ioproxy.WriterInstanceInterface
	Args    []string
	Utility utility.UtilityInterface
	Id      string
	Force   bool
	Verbose bool
}

func NewLikeArtistCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &likeArtistOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    g.Args,
		Utility: g.Utility,
	}

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.LIKE_ARTIST_USE
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.likeArtistRunE

	cmd.PersistentFlags().StringVarP(
		&o.Id,
		constant.LIKE_ARTIST_FLAG_ID,
		constant.LIKE_ARTIST_FLAG_ID_SHORTHAND,
		"",
		constant.LIKE_ARTIST_FLAG_ID_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Force,
		constant.LIKE_ARTIST_FLAG_FORCE,
		constant.LIKE_ARTIST_FLAG_FORCE_SHORTHAND,
		false,
		constant.LIKE_ARTIST_FLAG_FORCE_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Verbose,
		constant.LIKE_ARTIST_FLAG_VERBOSE,
		constant.LIKE_ARTIST_FLAG_VERBOSE_SHORTHAND,
		false,
		constant.LIKE_ARTIST_FLAG_VERBOSE_DESCRIPTION,
	)

	o.Out = g.Out
	o.ErrOut = g.ErrOut
	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	cmd.SetHelpTemplate(constant.LIKE_ARTIST_HELP_TEMPLATE)

	return cmd
}

func (o *likeArtistOption) likeArtistRunE(_ *cobra.Command, _ []string) error {
	o.Client = g.Client
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

	return o.likeArtist()
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
	colorProxy := colorproxy.New()
	if result.AlreadyLiked {
		return colorProxy.BlueString(fmt.Sprintf(like_artist_message_template_like_artist_already_liked, util.STRING_ARTIST, result.ArtistNames))
	} else if result.Refused {
		return colorProxy.YellowString(fmt.Sprintf(like_artist_message_template_like_artist_refused, util.STRING_ARTIST, result.ArtistNames))
	} else {
		return colorProxy.GreenString(fmt.Sprintf(like_artist_message_template_like_artist_succeeded, util.STRING_ARTIST, result.ArtistNames))
	}
}
