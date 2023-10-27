/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yanosea/spotlike/client"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

type SearchResult struct {
	ID     string
	Type   string
	Name   string
	Artist string
}

var (
	contentType string
	query       string
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "You can get the ID of some contents in Spotify.",
	Long: `You can get the ID of some contents in Spotify.

You can choose a content type below.
  * artist
  * album`,

	RunE: func(cmd *cobra.Command, args []string) error {

		// check type
		if strings.ToLower(contentType) != "artist" && strings.ToLower(contentType) != "album" {
			return errors.New("'type' must be 'artist' or 'album'")
		}

		// get client
		var spt *client.SpotifyClient
		var err error
		if spt, err = client.New(); err != nil {
			return err
		}

		// search type
		var searchType spotify.SearchType
		if strings.ToLower(contentType) == "artist" {
			searchType = spotify.SearchTypeArtist
		} else if strings.ToLower(contentType) == "album" {
			searchType = spotify.SearchTypeAlbum
		}

		// search
		var result *spotify.SearchResult
		if result, err = spt.Client.Search(spt.Context, query, searchType, spotify.Limit(1)); err != nil {
			return err
		}

		// set search results
		var searchResult SearchResult
		if result.Artists != nil {
			searchResult = SearchResult{
				ID:   result.Artists.Artists[0].ID.String(),
				Type: "Artist",
				Name: result.Artists.Artists[0].Name,
			}
		}

		if result.Albums != nil {
			searchResult = SearchResult{
				ID:     result.Albums.Albums[0].ID.String(),
				Type:   "Album",
				Name:   result.Albums.Albums[0].Name,
				Artist: result.Albums.Albums[0].Artists[0].Name,
			}
		}

		// output
		fmt.Printf("ID: %s\n", searchResult.ID)
		fmt.Printf("Type: %s\n", searchResult.Type)
		fmt.Printf("Name: %s\n", searchResult.Name)
		if searchResult.Artist != "" {
			fmt.Printf("Artist: %s\n", searchResult.Artist)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// flag 'type'
	getCmd.Flags().StringVarP(&contentType, "type", "t", "", "type of content")
	if err := getCmd.MarkFlagRequired("type"); err != nil {
		return
	}

	// flag 'query'
	getCmd.Flags().StringVarP(&query, "query", "q", "", "search query")
	if err := getCmd.MarkFlagRequired("query"); err != nil {
		return
	}
}
