package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yanosea/spotlike/api"
	"github.com/yanosea/spotlike/constants"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// Search performs a Spotify search based on the specified search type and query.
func Search(st string, query string) error {
	// define search type
	searchType, err := defineSearchType(st)
	if err != nil {
		return err
	}

	// execute search by query
	searchResult := api.SearchByQuery(Client, searchType, query)

	// print search result
	printSearchResult(searchResult)

	return nil
}

// defineSearchType converts the search type string to the corresponding spotify.SearchType.
func defineSearchType(searchType string) (spotify.SearchType, error) {
	// define search type
	searchType = strings.ToLower(searchType)
	if st, ok := constants.SearchTypeMap[searchType]; ok {
		return st, nil
	}

	return 0, errors.New(constants.SearchFlagTypeInvalidErrorMessage)
}

// printSearchResult prints the search result to the console.
func printSearchResult(searchResult *api.SearchResult) {
	if searchResult.Result {
		// search succeeded
		fmt.Println(formatResult(constants.Id, searchResult.ID))
		fmt.Println(formatResult(constants.Type, searchResult.Type))
		fmt.Println(formatResult(constants.Artist, searchResult.ArtistNames))
		if searchResult.Type == constants.Album {
			fmt.Println(formatResult(constants.Album, searchResult.AlbumName))
		}
		if searchResult.Type == constants.Track {
			fmt.Println(formatResult(constants.Album, searchResult.AlbumName))
			fmt.Println(formatResult(constants.Track, searchResult.TrackName))
		}
	} else {
		// search failed
		fmt.Printf(constants.SearchFailedErrorResultMessageFormat, searchResult.Type)
		fmt.Println(formatErrorResult(searchResult.Error))
	}
}
