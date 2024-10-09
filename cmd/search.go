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

type searchOption struct {
	Out        ioproxy.WriterInstanceInterface
	ErrOut     ioproxy.WriterInstanceInterface
	Args       []string
	Utility    utility.UtilityInterface
	Query      string
	Number     int
	SearchType string
}

func NewSearchCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &searchOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    g.Args,
		Utility: g.Utility,
	}

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.SEARCH_USE
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.searchRunE

	cmd.PersistentFlags().StringVarP(
		&o.Query,
		constant.SEARCH_FLAG_QUERY,
		constant.SEARCH_FLAG_QUERY_SHORTHAND,
		"",
		constant.SEARCH_FLAG_QUERY_DESCRIPTION,
	)
	cmd.PersistentFlags().IntVarP(
		&o.Number,
		constant.SEARCH_FLAG_NUMBER,
		constant.SEARCH_FLAG_NUMBER_SHORTHAND,
		1,
		constant.SEARCH_FLAG_NUMBER_DESCRIPTION,
	)
	cmd.PersistentFlags().StringVarP(
		&o.SearchType,
		constant.SEARCH_FLAG_TYPE,
		constant.SEARCH_FLAG_TYPE_SHORTHAND,
		"",
		constant.SEARCH_FLAG_TYPE_DESCRIPTION,
	)

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.SEARCH_HELP_TEMPLATE)

	return cmd
}

func (o *searchOption) searchRunE(_ *cobra.Command, _ []string) error {
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

	return o.search()
}

func (o *searchOption) search() error {
	// set the query
	q := strings.TrimSpace(o.combineAllArgs() + o.Query)
	if q == "" {
		// if the args or the flag query is empty, return error
		return errors.New(search_error_message_args_or_flag_query_required)
	}
	// set the search type
	st := o.defineSearchType(o.SearchType)
	if st == 0 {
		// if the search type is invalid,  return error
		return errors.New(search_error_message_flag_type_invalid)
	}
	// execute search
	searchResult, err := api.SearchByQuery(o.Client, q, o.Number, st)
	if err != nil {
		return err
	}
	// print the search result
	o.printSearchResult(searchResult)

	return nil
}

func (o *searchOption) combineAllArgs() string {
	var allArgs string
	for index, arg := range o.Args {
		allArgs += arg
		if index+1 != len(o.Args) {
			// if the arg is not last, add space
			allArgs += " "
		}
	}

	return allArgs
}

func (o *searchOption) defineSearchType(searchType string) spotify.SearchType {
	if o.SearchType == "" {
		// if the search type is not defined,  set all types
		return spotify.SearchTypeArtist | spotify.SearchTypeArtist | spotify.SearchTypeTrack
	}
	if st, ok := util.SEARCH_TYPE_MAP[strings.ToLower(o.SearchType)]; ok {
		// if the search type is defined and valid,  return the search type
		return st
	}

	// if the search type is invalid, return 0
	return 0
}

func (o *searchOption) printSearchResult(searchResultList []api.SearchResult) {
	for index, searchResult := range searchResultList {
		if index != 0 {
			// if the result was not first, add blank line above
			util.PrintWithWriterWithBlankLineAbove(o.Out, formatSearchResult("🆔", util.STRING_ID, searchResult.Id))
		} else {
			util.PrintlnWithWriter(o.Out, formatSearchResult("🆔", util.STRING_ID, searchResult.Id))
		}
		util.PrintlnWithWriter(o.Out, formatSearchResult("📦", util.STRING_TYPE, util.SEARCH_TYPE_MAP_REVERSED[searchResult.Type]))
		util.PrintlnWithWriter(o.Out, formatSearchResult("🎤", util.STRING_ARTIST, searchResult.ArtistNames))
		if searchResult.Type == spotify.SearchTypeAlbum || searchResult.Type == spotify.SearchTypeTrack {
			// if the search type is album or track, print the release date and the album name
			util.PrintlnWithWriter(o.Out, formatSearchResult("📅", util.STRING_RELEASE, searchResult.ReleaseDate))
			util.PrintlnWithWriter(o.Out, formatSearchResult("💿", util.STRING_ALBUM, searchResult.AlbumName))
		}
		if searchResult.Type == spotify.SearchTypeTrack {
			// if the search type is track, print the track name
			util.PrintlnWithWriter(o.Out, formatSearchResult("🎵", util.STRING_TRACK, searchResult.TrackName))
		}
	}
}

func formatSearchResult(icon string, topic string, detail string) string {
	return fmt.Sprintf("%s %s:\t[%s]", icon, addSpacesToLength(topic, 8), detail)
}

func addSpacesToLength(str string, length int) string {
	spacesToAdd := length - len(str)
	if spacesToAdd <= 0 {
		return str
	}
	return str + strings.Repeat(" ", spacesToAdd)
}
