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
)

var (
	// id : id of the content
	id string
	// unit : unit for like
	unit string
)

// likeCmd : cobra like command
var likeCmd = &cobra.Command{
	Use:   "like",
	Short: "You can like the contents in Spotify by ID.",
	Long: `You can like the contents in Spotify by ID.

You can like the content(s) by the units below.
  * artist (If you pass the artist ID, spotlike will like the artist.)
  * album  (If you pass the artist ID, spotlike will like all albums released by the artist.)
  * track  (If you pass the artist ID, spotlike will like all tracks released by the artist.)
           (If you pass the album ID, spotlike will like all tracks contained in the the album.)`,

	// RunE : like command
	RunE: func(cmd *cobra.Command, args []string) error {
		// get the client
		var spt *api.SpotifyClient
		if client, err := api.GetClient(); err != nil {
			return err
		} else {
			// set client
			spt = client
		}

		// execute search by id
		var searchResult *api.SearchResult
		sr := api.SearchById(spt, id)
		if !sr.Result {
			return sr.Error
		} else {
			// the content was not found
			searchResult = sr
		}

		// artist
		var likeResults []*api.LikeResult
		if searchResult.Type == "Artist" {
			// validate unit
			if strings.ToLower(unit) == "artist" || unit == "" {
				// like artist
				likeResults = api.LikeArtistById(spt, searchResult.ID)
			} else if strings.ToLower(unit) == "album" {
				// like all albums released by the artist
				likeResults = api.LikeAllAlbumsReleasedByArtistById(spt, searchResult.ID)
			} else if strings.ToLower(unit) == "track" {
				// like all tracks released by the artist
				likeResults = api.LikeAllTracksReleasedByArtistById(spt, searchResult.ID)
			} else {
				// wrong unit passed
				return errors.New("You passed the artist ID, so you have to pass 'artist' for the unit option. Or you should not specify the unit option.")
			}
		}

		// album
		if searchResult.Type == "Album" {
			// validate unit
			if strings.ToLower(unit) == "album" || unit == "" {
				// like album
				likeResults = api.LikeAlbumById(spt, searchResult.ID)
			} else if strings.ToLower(unit) != "track" {
				// like all tracks on the album
				likeResults = api.LikeAllTracksOnAlbumById(spt, searchResult.ID)
			} else {
				// wrong unit passed
				return errors.New("You passed the album ID, so you have to pass 'album' for the unit option. Or you should not specify the unit option.")
			}
		}

		// track
		if searchResult.Type == "Track" {
			// validate unit
			if strings.ToLower(unit) == "track" || unit == "" {
				// like track
				likeResults = api.LikeTrackById(spt, searchResult.ID)
			} else {
				// wrong unit passed
				return errors.New("You passed the track ID, so you have to pass 'track' for the unit option. Or you should not specify the unit option.")
			}
		}

		// output like result
		for _, result := range likeResults {
			if result.Result {
				fmt.Printf("Like %s:%s succeeded!\n", result.Type, result.Name)
			} else {
				fmt.Printf("Like %s:%s failed...\n%s\n", result.Type, result.Name, result.Error)
			}
		}

		return nil
	},
}

// init : executed before like command executed
func init() {
	rootCmd.AddCommand(likeCmd)

	// validate flag 'id'
	likeCmd.Flags().StringVarP(&id, "id", "i", "", "id of the content")
	if err := likeCmd.MarkFlagRequired("id"); err != nil {
		return
	}

	// validate flag 'unit'
	likeCmd.Flags().StringVarP(&unit, "unit", "u", "", "unit for like")
}
