/*
Package cmd defines and configures the sub commands for the spotlike CLI application using Cobra.
*/
package cmd

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

// variables
var (
	// searchType is the type of the content for search
	searchType string
	// query is used for search
	query string
)

// searchCmd is the Cobra search sub command of spotlike
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for the ID of content in Spotify.",
	Long: `Search for the ID of content in Spotify.

You can search for content using the type option below:
  * artist
  * album
  * track`,

	// RunE is the function to search
	RunE: func(cmd *cobra.Command, args []string) error {
		// validate search type
		if strings.ToLower(searchType) != "artist" && strings.ToLower(searchType) != "album" && strings.ToLower(searchType) != "track" {
			return errors.New("The option 'type' must be 'artist', 'album', or 'track'")
		}

		// get the Spotify client
		var spt *spotify.Client
		if client, err := api.GetClient(); err != nil {
			return err
		} else {
			// set the client
			spt = client
		}

		// define the search type
		var st spotify.SearchType
		if strings.ToLower(searchType) == "artist" {
			st = spotify.SearchTypeArtist
		} else if strings.ToLower(searchType) == "album" {
			st = spotify.SearchTypeAlbum
		} else if strings.ToLower(searchType) == "track" {
			st = spotify.SearchTypeTrack
		}

		// execute search by query
		searchResult := api.SearchByQuery(spt, st, query)

		// print search result
		if searchResult.Result {
			fmt.Printf("ID\t:\t%s\n", searchResult.ID)
			fmt.Printf("Type\t:\t%s\n", searchResult.Type)
			fmt.Printf("Artist\t:\t%s\n", searchResult.ArtistNames)
			if searchResult.Type == "Album" {
				fmt.Printf("Album\t:\t%s\n", searchResult.AlbumName)
			}
			if searchResult.Type == "Track" {
				fmt.Printf("Album\t:\t%s\n", searchResult.AlbumName)
				fmt.Printf("Track\t:\t%s\n", searchResult.TrackName)
			}
		} else {
			fmt.Printf("Search for %s failed...\n", searchResult.Type)
			fmt.Printf("query\t:\t%s\n", query)
			fmt.Printf("Error\t:\t%s\n", searchResult.Error)
		}
		return nil
	},
}

// init is executed before the search command is executed
func init() {
	rootCmd.AddCommand(searchCmd)

	// validate the flag 'type'
	searchCmd.Flags().StringVarP(&searchType, "type", "t", "", "Type of the content for search")
	if err := searchCmd.MarkFlagRequired("type"); err != nil {
		return
	}

	// validate the flag 'query'
	searchCmd.Flags().StringVarP(&query, "query", "q", "", "Query for search")
	if err := searchCmd.MarkFlagRequired("query"); err != nil {
		return
	}
}
