package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yanosea/spotlike/api"
	"github.com/yanosea/spotlike/util"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// constants
const (
	// search_use is the useage of search command.
	search_use = "search"
	// search_short is the short description of search command.
	search_short = "Search for the ID of content in Spotify."
	// search_long is the long description of search command.
	search_long = `Search for the ID of content in Spotify.

You can set the args or the flag "query" to specify the search query.
If you set both args and flag "query", they will be combined.

You can set the flag "type" to search type of the content.
If you don't set the flag "type", searching without specifying the content type will be executed.
You must specify the the flag "type" below :
  * artist
  * album
  * track`
	// search_flag_query is the string of the query flag of the search command.
	search_flag_query = "query"
	// search_shorhand_query is the string of the query shorthand flag of the search command.
	search_shorhand_query = "q"
	// search_flag_description_query the description of the query flag of the search command.
	search_flag_description_query = "query for search"
	// search_flag_type is the string of the type flag of the search command.
	search_flag_type = "type"
	// search_shorthand_type  is the string of the type shorthand flag of the search command.
	search_shorthand_type = "t"
	// search_flag_description_type  is the description of the type flag of the search command.
	search_flag_description_type = "type of the content for search"
	// search_error_message_args_or_query_required is the error message for the invalid arguments or the query.
	search_error_message_args_or_flag_query_required = `The arguments or the flag "query" is required...`
	// search_error_message_flag_type_invalid is the error message for the invalid type.
	search_error_message_flag_type_invalid = `The argument of the flag "type" must be "artist", "album", or "track"...`
)

// variables
var (
	// query holds flag 'query'
	query string
	// searchType holds flag 'type'
	searchType string
)

// searchCmd is the Cobra search sub command of spotlike.
var searchCmd = &cobra.Command{
	Use:   search_use,
	Short: search_short,
	Long:  search_long,

	// RunE is the function to search.
	RunE: func(cmd *cobra.Command, args []string) error {
		// execute search
		if err := search(args, query, searchType); err != nil {
			return err
		}

		return nil
	},
}

// init is executed before the search command is executed.
func init() {
	// cobra init
	rootCmd.AddCommand(searchCmd)
	// set the flag 'type'
	searchCmd.Flags().StringVarP(&searchType, search_flag_type, search_shorthand_type, "", search_flag_description_type)
	// set the flag 'query'
	searchCmd.Flags().StringVarP(&query, search_flag_query, search_shorhand_query, "", search_flag_description_query)
}

// search performs a Spotify search based on the specified search type and query.
func search(args []string, query string, searchType string) error {
	// combine all the args and flag "query" then check it
	q := strings.TrimSpace(combineAllArgs(args) + query)
	if q == "" {
		return errors.New(search_error_message_args_or_flag_query_required)
	}

	// define the search type and check it
	st, err := defineSearchType(searchType)
	if err != nil {
		return err
	}
	if st == 0 {
		// search without specifying type
		st = spotify.SearchTypeArtist | spotify.SearchTypeArtist | spotify.SearchTypeTrack
	}

	// execute search by the query and print the search result
	printSearchResult(api.SearchByQuery(Client, st, q))

	return nil
}

// combineAllArgs combines all the the args and returns it as string
func combineAllArgs(args []string) string {
	var allArgs string
	for index, arg := range args {
		allArgs += arg
		if index+1 != len(args) {
			allArgs += " "
		}
	}

	return allArgs
}

// defineSearchType converts the search type string to the corresponding spotify.SearchType.
func defineSearchType(searchType string) (spotify.SearchType, error) {
	// define search type
	if searchType == "" {
		return 0, nil

	}
	if st, ok := util.SEARCH_TYPE_MAP[strings.ToLower(searchType)]; ok {
		return st, nil
	}

	return 0, errors.New(search_error_message_flag_type_invalid)
}

// printSearchResult prints the search result to the console.
func printSearchResult(searchResult *api.SearchResult, err error) {
	if err != nil {
		// search failed
		fmt.Println(formatSearchResultError(err))

		return
	}

	// search succeeded
	fmt.Println(formatSearchResult(util.STRING_ID, searchResult.Id))
	fmt.Println(formatSearchResult(util.STRING_TYPE, util.SEARCH_TYPE_MAP_REVERSED[searchResult.Type]))
	fmt.Println(formatSearchResult(util.STRING_ARTIST, searchResult.ArtistNames))

	if searchResult.Type == spotify.SearchTypeAlbum {
		// search album
		fmt.Println(formatSearchResult(util.STRING_ALBUM, searchResult.AlbumName))
	}

	if searchResult.Type == spotify.SearchTypeTrack {
		// search track
		fmt.Println(formatSearchResult(util.STRING_ALBUM, searchResult.AlbumName))
		fmt.Println(formatSearchResult(util.STRING_TRACK, searchResult.TrackName))
	}
}

// formatSearchResult returns the formatted search result
func formatSearchResult(topic string, detail string) string {
	return fmt.Sprintf("%s\t:\t%s", topic, detail)
}

// formatSearchResultError returns the formatted search error result
func formatSearchResultError(error error) string {
	return fmt.Sprintf("Error:\n  %s", error)
}
