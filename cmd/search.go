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

You can search for content using the type option below:
  * artist
  * album
  * track`
	// search_flag_type is the string of the type flag of the search command.
	search_flag_type = "type"
	// search_shorthand_type  is the string of the type shorthand flag of the search command.
	search_shorthand_type = "t"
	// search_flag_description_type  is the description of the type flag of the search command.
	search_flag_description_type = "type of the content for search"
	// search_flag_query is the string of the query flag of the search command.
	search_flag_query = "query"
	// search_shorhand_query is the string of the query shorthand flag of the search command.
	search_shorhand_query = "q"
	// search_flag_description_query the description of the query flag of the search command.
	search_flag_description_query = "query for search"
	// search_error_message_flag_type_invalid is the error message for the invalid type.
	search_error_message_flag_type_invalid = `The argument of the flag "type" must be "artist", "album", or "track..."`
)

// variables
var (
	// searchType holds flag 'type'
	searchType string
	// query holds flag 'query'
	query string
)

// searchCmd is the Cobra search sub command of spotlike.
var searchCmd = &cobra.Command{
	Use:   search_use,
	Short: search_short,
	Long:  search_long,

	// RunE is the function to search.
	RunE: func(cmd *cobra.Command, args []string) error {
		// execute search
		if err := search(); err != nil {
			return err
		}

		return nil
	},
}

// init is executed before the search command is executed.
func init() {
	// cobra init
	rootCmd.AddCommand(searchCmd)

	// validate the flag 'type'
	searchCmd.Flags().StringVarP(&searchType, search_flag_type, search_shorthand_type, "", search_flag_description_type)
	if err := searchCmd.MarkFlagRequired(search_flag_type); err != nil {
		return
	}

	// validate the flag 'query'
	searchCmd.Flags().StringVarP(&query, search_flag_query, search_shorhand_query, "", search_flag_description_query)
	if err := searchCmd.MarkFlagRequired(search_flag_query); err != nil {
		return
	}
}

// search performs a Spotify search based on the specified search type and query.
func search() error {
	// define search type
	st, err := defineSearchType(searchType)
	if err != nil {
		return err
	}

	// execute search by query
	searchResult, err := api.SearchByQuery(Client, st, query)

	// print search result
	printSearchResult(searchResult, err)

	return nil
}

// defineSearchType converts the search type string to the corresponding spotify.SearchType.
func defineSearchType(searchType string) (spotify.SearchType, error) {
	// define search type
	searchType = strings.ToLower(searchType)
	if st, ok := util.SEARCH_TYPE_MAP[searchType]; ok {
		return st, nil
	}

	return 0, errors.New(search_error_message_flag_type_invalid)
}

// printSearchResult prints the search result to the console.
func printSearchResult(searchResult *api.SearchResult, err error) {
	if err == nil {
		// search succeeded
		fmt.Println(formatSearchResult(util.STRING_ID, searchResult.Id))
		fmt.Println(formatSearchResult(util.STRING_TYPE, util.SEARCH_TYPE_MAP_REVERSED[searchResult.Type]))
		fmt.Println(formatSearchResult(util.STRING_ARTIST, searchResult.ArtistNames))
		if searchResult.Type == spotify.SearchTypeAlbum {
			fmt.Println(formatSearchResult(util.STRING_ALBUM, searchResult.AlbumName))
		}
		if searchResult.Type == spotify.SearchTypeTrack {
			fmt.Println(formatSearchResult(util.STRING_ALBUM, searchResult.AlbumName))
			fmt.Println(formatSearchResult(util.STRING_TRACK, searchResult.TrackName))
		}
	} else {
		// search failed
		fmt.Println(formatSearchErrorResult(err))
	}
}

// formatSearchResult returns the formatted result
func formatSearchResult(topic string, detail string) string {
	return fmt.Sprintf("%s\t:\t%s", topic, detail)
}

// formatSearchErrorResult returns the formatted error result
func formatSearchErrorResult(error error) string {
	return fmt.Sprintf("Error:\n  %s", error)
}
