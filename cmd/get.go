/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
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
		fmt.Printf("Your mind number is %s !\n", contentType)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&contentType, "type", "t", "", "type of content")
	if err := getCmd.MarkFlagRequired("type"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	getCmd.Flags().StringVarP(&query, "query", "q", "", "search query")
	if err := getCmd.MarkFlagRequired("query"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
