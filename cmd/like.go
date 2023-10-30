/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	"fmt"

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
  * album (If you pass the artist ID, spotlike will like all albums released by the artist.)
  * track (If you pass the artist ID, spotlike will like all tracks released by the artist.)
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
		if searchResult, err := api.SearchById(spt, id); err != nil {
			return err
		} else {
			// output search result
			fmt.Printf("ID:\t%s\n", searchResult.ID)
			fmt.Printf("Type:\t%s\n", searchResult.Type)
			fmt.Printf("Name:\t%s\n", searchResult.Name)

			if searchResult.Album != "" {
				fmt.Printf("Album:\t%s\n", searchResult.Album)
			}

			if searchResult.Artist != "" {
				fmt.Printf("Artist:\t%s\n", searchResult.Artist)
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
