package cmd

import (
	"github.com/yanosea/spotlike/cli"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

// variables
var (
	// searchType is the type of the content for search
	searchType string
	// query is used for search
	query string
)

// searchCmd is the Cobra search sub command of spotlike
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for the ID of content in Spotify.",
	Long: `Search for the ID of content in Spotify.

You can search for content using the type option below:
  * artist
  * album
  * track`,

	// RunE is the function to search
	RunE: func(cmd *cobra.Command, args []string) error {
		// execute search
		if err := cli.Search(searchType, query); err != nil {
			return err
		}

		return nil
	},
}

// init is executed before the search command is executed
func init() {
	// cobra init
	rootCmd.AddCommand(searchCmd)

	// validate the flag 'type'
	searchCmd.Flags().StringVarP(&searchType, "type", "t", "", "type of the content for search")
	if err := searchCmd.MarkFlagRequired("type"); err != nil {
		return
	}

	// validate the flag 'query'
	searchCmd.Flags().StringVarP(&query, "query", "q", "", "query for search")
	if err := searchCmd.MarkFlagRequired("query"); err != nil {
		return
	}
}
