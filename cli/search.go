package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yanosea/spotlike/api"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// st is the type of the content for search
var st spotify.SearchType

func Search(searchType string, query string) error {
	// degine search type
	if err := defineSearchType(searchType); err != nil {
		return err
	}

	// execute search by query
	searchResult := api.SearchByQuery(Client, st, query)

	// print search result
	printSearchResult(searchResult)

	return nil
}

func defineSearchType(searchType string) error {
	// define search type
	switch strings.ToLower(searchType) {
	case "artist":
		st = spotify.SearchTypeArtist
	case "album":
		st = spotify.SearchTypeAlbum
	case "track":
		st = spotify.SearchTypeTrack
	default:
		return errors.New(`the argument of the flag "type" must be "artist", "album", or "track"`)
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
