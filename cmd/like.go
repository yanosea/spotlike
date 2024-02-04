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

You must set the args or the flag "id" to content(s) for like.
If you set both args and flag "id", both contents will be liked.

You can like content(s) using the flag "level" below :
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
	like_flag_description_force = "like without confirming"
	// search_error_message_args_or_query_required is the error message for the invalid arguments or the query.
	like_error_message_args_or_flag_id_required = `The arguments or the flag "id" is required...`
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
	// id holds flag 'id'
	id string
	// level holds flag 'level'
	level string
	// force holds flag 'force'
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
		if err := like(args, id, level, force); err != nil {
			return err
		}

		return nil
	},
}

// init is executed before the like command is executed
func init() {
	// cobra init
	rootCmd.AddCommand(likeCmd)
	// validate the flag 'id'
	likeCmd.Flags().StringVarP(&id, like_flag_id, like_shorthand_id, "", like_flag_description_id)
	// validate the flag 'level'
	likeCmd.Flags().StringVarP(&level, like_flag_level, like_shorthand_level, "", like_flag_description_level)
	// validate the flag 'force'
	likeCmd.Flags().BoolVarP(&force, like_flag_force, like_shorthand_force, false, like_flag_description_force)
}

// like performs a Spotify like based on the specified id and level.
func like(args []string, id string, level string, force bool) error {
	// separete all the args and the flag "id" then check it
	var ids []string
	for _, arg := range args {
		a := strings.Fields(arg)
		for _, b := range a {
			ids = append(ids, b)
		}
	}
	a := strings.Fields(id)
	for _, b := range a {
		ids = append(ids, b)
	}
	if len(ids) == 0 {
		return errors.New(like_error_message_args_or_flag_id_required)
	}

	for _, id := range ids {
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
			} else if strings.ToLower(level) == util.SEARCH_TYPE_MAP_REVERSED[spotify.SearchTypeTrack] {
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
	}

	return nil
}

// printLikeResult prints the like result to the console.
func printLikeResult(likeResults []*api.LikeResult) {
	// for each result
	for _, result := range likeResults {
		// like failed
		if result.Error != nil {
			fmt.Println(formatLikeResultError(result.Error))
		}

		if result.Type == spotify.SearchTypeArtist {
			// like artist
			fmt.Println(formatLikeArtistResult(result))
		}

		if result.Type == spotify.SearchTypeAlbum {
			// like album
			fmt.Println(formatLikeAlbumResult(result))
		}

		if result.Type == spotify.SearchTypeTrack {
			// like track
			fmt.Println(formatLikeTrackResult(result))
		}
	}
}

// formatLikeArtistResult returns the formatted like artist result
func formatLikeArtistResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		// already liked
		return fmt.Sprintf("%s already liked!\t:\t[%s]", util.STRING_ARTIST, result.ArtistNames)
	} else if result.Refused {
		// refused
		return fmt.Sprintf("Like %s refused!\t:\t[%s]", util.STRING_ARTIST, result.ArtistNames)
	} else {
		// liked
		return fmt.Sprintf("Like %s succeeded!\t:\t[%s]", util.STRING_ARTIST, result.ArtistNames)
	}
}

// formatLikeAlbumResult returns the formatted like album result
func formatLikeAlbumResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		// already liked
		return fmt.Sprintf("%s by %s already liked!\t:\t[%s]", util.STRING_ALBUM, result.ArtistNames, result.AlbumName)
	} else if result.Refused {
		// refused
		return fmt.Sprintf("Like %s by %s refused!\t:\t[%s]", util.STRING_ALBUM, result.ArtistNames, result.AlbumName)
	} else {
		// liked
		return fmt.Sprintf("Like %s by %s succeeded!\t:\t[%s]", util.STRING_ALBUM, result.ArtistNames, result.AlbumName)
	}
}

// formatLikeTrackResult returns the formatted like track result
func formatLikeTrackResult(result *api.LikeResult) string {
	if result.AlreadyLiked {
		// already liked
		return fmt.Sprintf("%s in %s by %s already liked!\t:\t[%s]", util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName)
	} else if result.Refused {
		// refused
		return fmt.Sprintf("Like %s in %s by %s refused!\t:\t[%s]", util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName)
	} else {
		// liked
		return fmt.Sprintf("Like %s in %s by %s succeeded!\t:\t[%s]", util.STRING_TRACK, result.AlbumName, result.ArtistNames, result.TrackName)
	}
}

// formatLikeResultError returns the formatted like error result
func formatLikeResultError(error error) string {
	return fmt.Sprintf("Error:\n  %s", error)
}
