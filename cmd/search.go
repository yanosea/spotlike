/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
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

var (
	// searchType : type of the content for search
	searchType string
	// query : query for search
	query string
)

// searchCmd : cobra search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "You can search the ID of the content in Spotify.",
	Long: `You can search the ID of the content in Spotify.

You can search the content using type option below.
  * artist
  * album
  * track`,

	// RunE : search command
	RunE: func(cmd *cobra.Command, args []string) error {
		// validate search type
		if strings.ToLower(searchType) != "artist" && strings.ToLower(searchType) != "album" && strings.ToLower(searchType) != "track" {
			return errors.New("'type' must be 'artist', 'album' or 'track'")
		}

		// get the client
		var spt *spotify.Client
		if client, err := api.GetClient(); err != nil {
			return err
		} else {
			// set client
			spt = client
		}

		// define search type
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

		if searchResult.Result {
			// output search result
			fmt.Printf("ID\t:\t%s\n", searchResult.ID)
			fmt.Printf("Type\t:\t%s\n", searchResult.Type)
			fmt.Printf("Name\t:\t%s\n", searchResult.Name)
			if searchResult.Album != "" {
				fmt.Printf("Album\t:\t%s\n", searchResult.Album)
			}
			if searchResult.Artist != "" {
				fmt.Printf("Artist\t:\t%s\n", searchResult.Artist)
			}
		} else {
			fmt.Printf("Search %s failed...\n", searchResult.Type)
			fmt.Printf("query\t:\t%s\n", query)
			fmt.Printf("Error\t:\t%s\n", searchResult.Error)
		}

		return nil
	},
}

// init : executed before search command executed
func init() {
	rootCmd.AddCommand(searchCmd)

	// validate the flag 'type'
	searchCmd.Flags().StringVarP(&searchType, "type", "t", "", "type of the content for search")
	if err := searchCmd.MarkFlagRequired("type"); err != nil {
		return
	}

	// validate the flag 'query'
	searchCmd.Flags().StringVarP(&query, "query", "q", "", "query for search")
	if err := searchCmd.MarkFlagRequired("query"); err != nil {
		return
	}
}
