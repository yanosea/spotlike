package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yanosea/spotlike/api"
	"github.com/yanosea/spotlike/util"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// constants
const (
	like_use   = "like"
	like_short = "Like content in Spotify by ID."
	like_long  = `Like content in Spotify by ID.

You can like content(s) using the level option below:
  * artist (If you pass the artist ID, spotlike will like the artist.)
  * album  (If you pass the artist ID, spotlike will like all albums released by the artist.)
           (If you pass the album ID, spotlike will like the album.)
  * track  (If you pass the artist ID, spotlike will like all tracks released by the artist.)
           (If you pass the album ID, spotlike will like all tracks contained in the album.)
           (If you pass the track ID, spotlike will like the track.)`
	// like_flag_id is the string of the id flag of the like command.
	like_flag_id = "id"
	// like_shorthand_id is the string of the id shorthand flag of the like command.
	like_shorthand_id = "i"
	// search_flag_description_type  is the description of the type flag of the like command.
	like_flag_description_id = "ID of the content for like"
	// like_flag_level is the string of the level flag of the like command.
	like_flag_level = "level"
	// like_shorthand_level is the string of the level shorthand flag of the like command.
	like_shorthand_level = "l"
	// search_flag_description_level is the description of the level flag of the like command.
	like_flag_description_level = "level for like"
	// like_flag_force is the string of the force flag of the like command.
	like_flag_force = "force"
	// like_shorthand_force is the string of the force shorthand flag of the like command.
	like_shorthand_force = "f"
	// search_flag_description_force is the description of the force flag of the like command.
	like_flag_description_force = "whether to like even if the content heve already been liked"
	// like_error_message_flag_level_invalid_artist is the error message for the invalid level with artist id.
	like_error_message_flag_level_invalid_artist = "You passed the artist ID, so you have to pass 'artist', 'album', or 'track' for the level option. Or you should not specify the level option to like the artist."
	// like_error_message_flag_level_invalid_album is the error message for the invalid level with album id.
	like_error_message_flag_level_invalid_album = "You passed the album ID, so you have to pass 'album' or 'track' for the level option. Or you should not specify the level option to like the album."
	// like_error_message_flag_level_invalid_track is the error message for the invalid level with track id.
	like_error_message_flag_level_invalid_track = "You passed the track ID, so you have to pass 'track' for the level option. Or you should not specify the level option to like the track."
	// like_error_message_something_wrong is the error message for something wrong searching before liking.
	like_error_message_something_wrong_with_searching = "Search result is wrong..."
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
	Use:   like_use,
	Short: like_short,
	Long:  like_long,

	// RunE is the function to like
	RunE: func(cmd *cobra.Command, args []string) error {
		// execute like
		if err := like(id, level, force); err != nil {
			return err
		}

		return nil
	},
}

// init is executed before the like command is executed
func init() {
	// cobra init
	rootCmd.AddCommand(likeCmd)

	// validate flag 'id'
	likeCmd.Flags().StringVarP(&id, like_flag_id, like_shorthand_id, "", like_flag_description_id)
	if err := likeCmd.MarkFlagRequired(like_flag_id); err != nil {
		return
	}

	// validate flag 'level'
	likeCmd.Flags().StringVarP(&level, like_flag_level, like_shorthand_level, "", like_flag_description_level)

	// validate flag 'force'
	likeCmd.Flags().BoolVarP(&force, like_flag_force, like_shorthand_force, false, like_flag_description_force)
}

// like performs a Spotify like based on the specified id and level.
func like(id string, level string, force bool) error {
	// first, execute search by ID
	searchResult, err := api.SearchById(Client, id)
	if err != nil {
		// the content was not found
		return err
	}

	// execute like
	var likeResults []*api.LikeResult
	switch searchResult.Type {
	// the type of the content is artist
	case spotify.SearchTypeArtist:
		// validate level
		if strings.ToLower(level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeArtist] || level == "" {
			// like artist
			likeResults = api.LikeArtistById(Client, searchResult.Id, force)
		} else if strings.ToLower(level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeAlbum] {
			// like all albums released by the artist
			likeResults = api.LikeAllAlbumsReleasedByArtistById(Client, searchResult.Id, force)
		} else if strings.ToLower(level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeAlbum] {
			// like all tracks released by the artist
			likeResults = api.LikeAllTracksReleasedByArtistById(Client, searchResult.Id, force)
		} else {
			// wrong level passed
			return errors.New(like_error_message_flag_level_invalid_artist)
		}

	// the type of the content is album
	case spotify.SearchTypeAlbum:
		// validate level
		if strings.ToLower(level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeAlbum] || level == "" {
			// like album
			likeResults = api.LikeAlbumById(Client, searchResult.Id, force)
		} else if strings.ToLower(level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeTrack] {
			// like all tracks on the album
			likeResults = api.LikeAllTracksInAlbumById(Client, searchResult.Id, force)
		} else {
			// wrong level passed
			return errors.New(like_error_message_flag_level_invalid_album)
		}

	// the type of the content is track
	case spotify.SearchTypeTrack:
		// validate level
		if strings.ToLower(level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeTrack] || level == "" {
			// like track
			likeResults = api.LikeTrackById(Client, searchResult.Id, force)
		} else {
			// wrong level passed
			return errors.New(like_error_message_flag_level_invalid_track)
		}

	default:
		return errors.New(like_error_message_something_wrong_with_searching)
	}

	// print like result
	printLikeResult(likeResults)

	return nil
}

// printLikeResult prints the like result to the console.
func printLikeResult(likeResults []*api.LikeResult) {
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
}
