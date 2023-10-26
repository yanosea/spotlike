/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/yanosea/spotlike/client"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

var (
	Client      *spotify.Client
	contentType string
	query       string
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "You can get the ID of some contents in Spotify.",
	Long: `You can get the ID of some contents in Spotify.

You can choose a content type below.
	* album
	* artist`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	if client, err := client.New(); err != nil {
		fmt.Println(err)
		return
	} else {
		Client = client
	}

	getCmd.Flags().StringVarP(&contentType, "type", "t", "", "type of content")
	if err := getCmd.MarkFlagRequired("type"); err != nil {
		fmt.Println(err)
		return
	}
	getCmd.Flags().StringVarP(&query, "query", "q", "", "search query")
	if err := getCmd.MarkFlagRequired("query"); err != nil {
		fmt.Println(err)
		return
	}
}
