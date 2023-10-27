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
	Album  string
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
  * album
	* track`,

	RunE: func(cmd *cobra.Command, args []string) error {

		// check type
		if strings.ToLower(contentType) != "artist" && strings.ToLower(contentType) != "album" && strings.ToLower(contentType) != "track" {
			return errors.New("'type' must be 'artist', 'album' or 'track'")
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
		} else if strings.ToLower(contentType) == "track" {
			searchType = spotify.SearchTypeTrack
		}

		// search
		var result *spotify.SearchResult
		if result, err = spt.Client.Search(spt.Context, query, searchType, spotify.Limit(1)); err != nil {
			return err
		}

		// set search results
		var searchResult SearchResult

		// artist
		if result.Artists != nil {
			searchResult = SearchResult{
				ID:   result.Artists.Artists[0].ID.String(),
				Type: "Artist",
				Name: result.Artists.Artists[0].Name,
			}
		}

		// album
		if result.Albums != nil {
			searchResult = SearchResult{
				ID:     result.Albums.Albums[0].ID.String(),
				Type:   "Album",
				Name:   result.Albums.Albums[0].Name,
				Artist: result.Albums.Albums[0].Artists[0].Name,
			}
		}

		// track
		if result.Tracks != nil {
			searchResult = SearchResult{
				ID:     result.Tracks.Tracks[0].ID.String(),
				Type:   "Track",
				Name:   result.Tracks.Tracks[0].Name,
				Album:  result.Tracks.Tracks[0].Album.Name,
				Artist: result.Tracks.Tracks[0].Artists[0].Name,
			}
		}

		// output
		fmt.Printf("ID: %s\n", searchResult.ID)
		fmt.Printf("Type: %s\n", searchResult.Type)
		fmt.Printf("Name: %s\n", searchResult.Name)
		if searchResult.Album != "" {
			fmt.Printf("Album: %s\n", searchResult.Album)
		}
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
