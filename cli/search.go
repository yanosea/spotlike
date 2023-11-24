package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yanosea/spotlike/api"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// st is the type of the content for search
var st spotify.SearchType

func Search(searchType string, query string) error {
	// execute search by query
	searchResult := api.SearchByQuery(Client, st, query)

	// print search result
	printSearchResult(searchResult)

	return nil
}

func ValidateFlags(searchCmd *cobra.Command, searchType string, query string) error {
	// validate the flag 'type'
	if err := searchCmd.MarkFlagRequired("type"); err != nil {
		return err
	}
	if strings.ToLower(searchType) != "artist" && strings.ToLower(searchType) != "album" && strings.ToLower(searchType) != "track" {
		return errors.New("The option 'type' must be 'artist', 'album', or 'track'")
	} else {
		// define search type
		if strings.ToLower(searchType) == "artist" {
			st = spotify.SearchTypeArtist
		} else if strings.ToLower(searchType) == "album" {
			st = spotify.SearchTypeAlbum
		} else if strings.ToLower(searchType) == "track" {
			st = spotify.SearchTypeTrack
		}
	}

	// validate the flag 'query'
	if err := searchCmd.MarkFlagRequired("query"); err != nil {
		return err
	}

	return nil
}

func printSearchResult(searchResult *api.SearchResult) {
	if searchResult.Result {
		fmt.Println(formatResult("ID", searchResult.ID))
		fmt.Println(formatResult("Type", searchResult.Type))
		fmt.Println(formatResult("Artist", searchResult.ArtistNames))
		if searchResult.Type == "Album" {
			fmt.Println(formatResult("Album", searchResult.AlbumName))
		}
		if searchResult.Type == "Track" {
			fmt.Println(formatResult("Album", searchResult.AlbumName))
			fmt.Println(formatResult("Track", searchResult.TrackName))
		}
	} else {
		fmt.Printf("Search for %s failed...\n", searchResult.Type)
		fmt.Println(formatErrorResult(searchResult.Error))
	}
}
