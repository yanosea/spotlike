package cmd

import (
	"github.com/yanosea/spotlike/cli"
	"github.com/yanosea/spotlike/constants"

	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

// variables
var (
	// searchType is the type of the content for search.
	searchType string
	// query is used for search.
	query string
)

// searchCmd is the Cobra search sub command of spotlike.
var searchCmd = &cobra.Command{
	Use:   constants.SearchUse,
	Short: constants.SearchShort,
	Long:  constants.SearchLong,

	// RunE is the function to search.
	RunE: func(cmd *cobra.Command, args []string) error {
		// execute search
		if err := cli.Search(searchType, query); err != nil {
			return err
		}

		return nil
	},
}

// init is executed before the search command is executed.
func init() {
	// cobra init
	rootCmd.AddCommand(searchCmd)

	// validate the flag 'type'
	searchCmd.Flags().StringVarP(&searchType, constants.SearchFlagType, constants.SearchFlagTypeShorthand, "", constants.SearchFlagTypeDescription)
	if err := searchCmd.MarkFlagRequired(constants.SearchFlagType); err != nil {
		return
	}

	// validate the flag 'query'
	searchCmd.Flags().StringVarP(&query, constants.SearchFlagQuery, constants.SearchFlagQueryShorthand, "", constants.SearchFlagQueryDescription)
	if err := searchCmd.MarkFlagRequired(constants.SearchFlagQuery); err != nil {
		return
	}
}
