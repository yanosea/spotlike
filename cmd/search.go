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
	search_short = "🔍 Search for the ID of content in Spotify."
	search_long  = `🔍 Search for the ID of content in Spotify.

You can set the args or the flag "query" to specify the search query.
If you set both args and flag "query", they will be combined.

You can set the flag "number" to limiting the number of search results.
Default is 1.

You can set the flag "type" to search type of the content.
If you don't set the flag "type", searching without specifying the content type will be executed.
You must specify the the flag "type" below :

* 🖌️ artist
* 💿 album
* 🎵 track`
	search_flag_query                                = "query"
	search_flag_query_shorthand                      = "q"
	search_flag_query_description                    = "query for search"
	search_flag_number                               = "number"
	search_flag_number_shorthand                     = "n"
	search_flag_number_description                   = "number of search results"
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
	Number     int
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

			return o.search()
		},
	}

	cmd.PersistentFlags().StringVarP(&o.Query, search_flag_query, search_flag_query_shorthand, "", search_flag_query_description)
	cmd.PersistentFlags().IntVarP(&o.Number, search_flag_number, search_flag_number_shorthand, 1, search_flag_number_description)
	cmd.PersistentFlags().StringVarP(&o.SearchType, search_flag_type, search_flag_type_shorthand, "", search_flag_type_description)
	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	return cmd
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
			util.PrintWithWriterWithBlankLineAbove(o.Out, formatSearchResult(util.STRING_ID, searchResult.Id))
		} else {
			util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_ID, searchResult.Id))
		}
		util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_TYPE, util.SEARCH_TYPE_MAP_REVERSED[searchResult.Type]))
		util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_ARTIST, searchResult.ArtistNames))
		if searchResult.Type == spotify.SearchTypeAlbum || searchResult.Type == spotify.SearchTypeTrack {
			// if the search type is album or track, print the release date and the album name
			util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_RELEASE, searchResult.ReleaseDate))
			util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_ALBUM, searchResult.AlbumName))
		}
		if searchResult.Type == spotify.SearchTypeTrack {
			// if the search type is track, print the track name
			util.PrintlnWithWriter(o.Out, formatSearchResult(util.STRING_TRACK, searchResult.TrackName))
		}
	}
}

func formatSearchResult(topic string, detail string) string {
	return fmt.Sprintf("%s\t:\t[%s]", topic, detail)
}
