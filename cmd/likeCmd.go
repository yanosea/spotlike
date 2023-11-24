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
	// id is the ID of content
	id string
	// level is level for like
	level string
	// foece is the option whether to like even if the content heve already been liked
	force bool
)

// likeCmd is the Cobra like sub command of spotlike
var likeCmd = &cobra.Command{
	Use:   "like",
	Short: "Like content in Spotify by ID.",
	Long: `Like content in Spotify by ID.

You can like content(s) using the level option below:
  * artist (If you pass the artist ID, spotlike will like the artist.)
  * album  (If you pass the artist ID, spotlike will like all albums released by the artist.)
           (If you pass the album ID, spotlike will like the album.)
  * track  (If you pass the artist ID, spotlike will like all tracks released by the artist.)
           (If you pass the album ID, spotlike will like all tracks contained in the album.)
           (If you pass the track ID, spotlike will like the track.)`,

	// RunE is the function to like
	RunE: func(cmd *cobra.Command, args []string) error {
		// get the Spotify client
		var spt *spotify.Client
		if client, err := api.GetClient(); err != nil {
			return err
		} else {
			// set the client
			spt = client
		}

		// first, execute search by ID
		var searchResult *api.SearchResult
		sr := api.SearchById(spt, id)
		if !sr.Result {
			// the content was not found
			return sr.Error
		} else {
			searchResult = sr
		}

		// execute like
		var likeResults []*api.LikeResult
		switch searchResult.Type {
		// the type of the content is artist
		case "Artist":
			// validate level
			if strings.ToLower(level) == "artist" || level == "" {
				// like artist
				likeResults = api.LikeArtistById(spt, searchResult.ID, force)
			} else if strings.ToLower(level) == "album" {
				// like all albums released by the artist
				likeResults = api.LikeAllAlbumsReleasedByArtistById(spt, searchResult.ID, force)
			} else if strings.ToLower(level) == "track" {
				// like all tracks released by the artist
				likeResults = api.LikeAllTracksReleasedByArtistById(spt, searchResult.ID, force)
			} else {
				// wrong level passed
				return errors.New("\n  You passed the artist ID, so you have to pass 'artist', 'album', or 'track' for the level option. Or you should not specify the level option to like the artist.\n")
			}

		// the type of the content is album
		case "Album":
			// validate level
			if strings.ToLower(level) == "album" || level == "" {
				// like album
				likeResults = api.LikeAlbumById(spt, searchResult.ID, force)
			} else if strings.ToLower(level) == "track" {
				// like all tracks on the album
				likeResults = api.LikeAllTracksInAlbumById(spt, searchResult.ID, force)
			} else {
				// wrong level passed
				return errors.New("\n  You passed the album ID, so you have to pass 'album' or 'track' for the level option. Or you should not specify the level option to like the album.\n")
			}

		// the type of the content is track
		case "Track":
			// validate level
			if strings.ToLower(level) == "track" || level == "" {
				// like track
				likeResults = api.LikeTrackById(spt, searchResult.ID, force)
			} else {
				// wrong level passed
				return errors.New("\n  You passed the track ID, so you have to pass 'track' for the level option. Or you should not specify the level option to like the track.\n")
			}

		default:
			return errors.New("\n  Search result is wrong.\n")
		}

		// print like result
		for _, result := range likeResults {
			if result.Result {
				if result.Type == "Artist" {
					if result.Skip {
						fmt.Printf("Like %s skipped!\t:\t[%s]\n", result.Type, result.ArtistNames)
					} else {
						fmt.Printf("Like %s succeeded!\t:\t[%s]\n", result.Type, result.ArtistNames)
					}
				} else if result.Type == "Album" {
					if result.Skip {
						fmt.Printf("Like %s by %s skipped!\t:\t[%s]\n", result.Type, result.ArtistNames, result.AlbumName)
					} else {
						fmt.Printf("Like %s by %s succeeded!\t:\t[%s]\n", result.Type, result.ArtistNames, result.AlbumName)
					}
				} else if result.Type == "Track" {
					if result.Skip {
						fmt.Printf("Like %s in %s by %s skipped!\t:\t[%s]\n", result.Type, result.AlbumName, result.ArtistNames, result.TrackName)
					} else {
						fmt.Printf("Like %s in %s by %s succeeded!\t:\t[%s]\n", result.Type, result.AlbumName, result.ArtistNames, result.TrackName)
					}
				}
			} else {
				if result.ErrorMessage == "" {
					if result.Type == "Artist" {
						fmt.Printf("Like %s failed...\t:\t[%s]\n", result.Type, result.ArtistNames)
					} else if result.Type == "Album" {
						fmt.Printf("Like %s by %s failed...\t:\t[%s]\n", result.Type, result.ArtistNames, result.AlbumName)
					} else if result.Type == "Track" {
						fmt.Printf("Like %s in %s by %s failed...\t:\t[%s]\n", result.Type, result.AlbumName, result.ArtistNames, result.TrackName)
					}
					fmt.Printf("Error:\n  %s\n\n", result.Error)
				} else {
					fmt.Printf("%s\n\n", result.ErrorMessage)
					fmt.Printf("Error:\n  %s\n", result.Error)
				}
			}
		}
		return nil
	},
}

// init is executed before the like command is executed
func init() {
	rootCmd.AddCommand(likeCmd)

	// validate flag 'id'
	likeCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the content for like")
	if err := likeCmd.MarkFlagRequired("id"); err != nil {
		return
	}

	// validate flag 'level'
	likeCmd.Flags().StringVarP(&level, "level", "l", "", "Level for like")

	// validate flag 'force'
	likeCmd.Flags().BoolVarP(&force, "force", "f", false, "whether to like even if the content heve already been liked")
}
