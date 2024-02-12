package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/yanosea/spotlike/api"
	"github.com/yanosea/spotlike/auth"
	"github.com/yanosea/spotlike/util"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

const (
	search_use   = "search"
	search_short = "Search for the ID of content in Spotify."
	search_long  = `Search for the ID of content in Spotify.

You can set the args or the flag "query" to specify the search query.
If you set both args and flag "query", they will be combined.

You can set the flag "type" to search type of the content.
If you don't set the flag "type", searching without specifying the content type will be executed.
You must specify the the flag "type" below :
  * artist
  * album
  * track`
	search_flag_query                                = "query"
	search_flag_query_shorthand                      = "q"
	search_flag_query_description                    = "query for search"
	search_flag_type                                 = "type"
	search_flag_type_shorthand                       = "t"
	search_flag_type_description                     = "type of the content for search"
	search_error_message_args_or_flag_query_required = `The arguments or the flag "query" is required...`
	search_error_message_flag_type_invalid           = `The argument of the flag "type" must be "artist", "album", or "track"...`
)

type searchOption struct {
	Client *spotify.Client

	Args       []string
	Query      string
	SearchType string

	Out    io.Writer
	ErrOut io.Writer
}

func newSearchCommand(globalOption *GlobalOption) *cobra.Command {
	o := &searchOption{}

	cmd := &cobra.Command{
		Use:   search_use,
		Short: search_short,
		Long:  search_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Client = globalOption.Client
			if o.Client == nil {
				client, err := auth.GetClient()
				if err != nil {
					return err
				}
				o.Client = client
			}
			o.Args = args
			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut

			return o.search()
		},
	}

	cmd.PersistentFlags().StringVarP(&o.SearchType, search_flag_type, search_flag_type_shorthand, "", search_flag_type_description)
	cmd.PersistentFlags().StringVarP(&o.Query, search_flag_query, search_flag_query_shorthand, "", search_flag_query_description)
	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
}

func (o *searchOption) search() error {
	// set the query
	q := strings.TrimSpace(o.combineAllArgs() + o.Query)
	if q == "" {
		return errors.New(search_error_message_args_or_flag_query_required)
	}
	// set the search type
	st := o.defineSearchType(o.SearchType)
	if st == 0 {
		return errors.New(search_error_message_flag_type_invalid)
	}
	// execute search
	searchResult, err := api.SearchByQuery(o.Client, st, q)
	if err != nil {
		return err
	}
	// show the result
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
		// if the search type is defined and matched,  set the type
		return st
	}

	return 0
}

func (o *searchOption) printSearchResult(searchResult *api.SearchResult) {
	util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_ID, searchResult.Id))
	util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_TYPE, util.SEARCH_TYPE_MAP_REVERSED[searchResult.Type]))
	util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_ARTIST, searchResult.ArtistNames))
	if searchResult.Type == spotify.SearchTypeAlbum || searchResult.Type == spotify.SearchTypeTrack {
		// if the search type is album or track, print the album name
		util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_ALBUM, searchResult.AlbumName))
	}
	if searchResult.Type == spotify.SearchTypeTrack {
		// if the search type is track, print the track name
		util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_TRACK, searchResult.TrackName))
	}
}

func formatSearchResult(topic string, detail string) string {
	return fmt.Sprintf("%s\t:\t[%s]", topic, detail)
}
