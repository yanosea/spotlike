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
	// id : id of the content
	id string
	// level : level for like
	level string
)

// likeCmd : cobra like command
var likeCmd = &cobra.Command{
	Use:   "like",
	Short: "You can like the contents in Spotify by ID.",
	Long: `You can like the contents in Spotify by ID.

You can like the content(s) using level option below.
  * artist (If you pass the artist ID, spotlike will like the artist.)
  * album  (If you pass the artist ID, spotlike will like all albums released by the artist.)
           (If you pass the album ID, spotlike will like the the album.)
  * track  (If you pass the artist ID, spotlike will like all tracks released by the artist.)
           (If you pass the album ID, spotlike will like all tracks contained in the the album.)
           (If you pass the track ID, spotlike will like the track.)`,

	// RunE : like command
	RunE: func(cmd *cobra.Command, args []string) error {
		// get the client
		var spt *spotify.Client
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
			// validate level
			if strings.ToLower(level) == "artist" || level == "" {
				// like artist
				likeResults = api.LikeArtistById(spt, searchResult.ID)
			} else if strings.ToLower(level) == "album" {
				// like all albums released by the artist
				likeResults = api.LikeAllAlbumsReleasedByArtistById(spt, searchResult.ID)
			} else if strings.ToLower(level) == "track" {
				// like all tracks released by the artist
				likeResults = api.LikeAllTracksReleasedByArtistById(spt, searchResult.ID)
			} else {
				// wrong level passed
				return errors.New("You passed the artist ID, so you have to pass 'artist', 'album', or 'track' for the level option. Or you should not specify the level option to like the artist.")
			}
		}

		// album
		if searchResult.Type == "Album" {
			// validate level
			if strings.ToLower(level) == "album" || level == "" {
				// like album
				likeResults = api.LikeAlbumById(spt, searchResult.ID)
			} else if strings.ToLower(level) == "track" {
				// like all tracks on the album
				likeResults = api.LikeAllTracksInAlbumById(spt, searchResult.ID)
			} else {
				// wrong level passed
				return errors.New("You passed the album ID, so you have to pass 'album' or 'track' for the level option. Or you should not specify the level option to like the album.")
			}
		}

		// track
		if searchResult.Type == "Track" {
			// validate level
			if strings.ToLower(level) == "track" || level == "" {
				// like track
				likeResults = api.LikeTrackById(spt, searchResult.ID)
			} else {
				// wrong level passed
				return errors.New("You passed the track ID, so you have to pass 'track' for the level option. Or you should not specify the level option to like the track.")
			}
		}

		// output like result
		for _, result := range likeResults {
			if result.Result {
				if result.Type == "artist" {
					fmt.Printf("Like %s succeeded!\t:\t[%s]\n", result.Type, result.ArtistName)
				} else if result.Type == "Album" {
					fmt.Printf("Like %s by %s succeeded!\t:\t[%s]\n", result.Type, result.ArtistName, result.AlbumName)
				} else if result.Type == "Track" {
					fmt.Printf("Like %s in %s by %s succeeded!\t:\t[%s]\n", result.Type, result.AlbumName, result.ArtistName, result.TrackName)
				}
			} else {
				if result.ErrorMessage == "" {
					if result.Type == "artist" {
						fmt.Printf("Like %s failed...\t:\t[%s]\n", result.Type, result.ArtistName)
					} else if result.Type == "Album" {
						fmt.Printf("Like %s by %s failed...\t:\t[%s]\n", result.Type, result.ArtistName, result.AlbumName)
					} else if result.Type == "Track" {
						fmt.Printf("Like %s in %s by %s failed...\t:\t[%s]\n", result.Type, result.AlbumName, result.ArtistName, result.TrackName)
					}
					fmt.Printf("Error :\t%s\n", result.Error)
				} else {
					fmt.Printf("%s\n", result.ErrorMessage)
					fmt.Printf("Error :\t%s\n", result.Error)
				}
			}
		}

		return nil
	},
}

// init : executed before like command executed
func init() {
	rootCmd.AddCommand(likeCmd)

	// validate flag 'id'
	likeCmd.Flags().StringVarP(&id, "id", "i", "", "id of the content for like")
	if err := likeCmd.MarkFlagRequired("id"); err != nil {
		return
	}

	// validate flag 'level'
	likeCmd.Flags().StringVarP(&level, "level", "l", "", "level for like")
}
