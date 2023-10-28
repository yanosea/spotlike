/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package cmd

import (
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

// likeCmd : cobra like command
var likeCmd = &cobra.Command{
	Use:   "like",
	Short: "You can like the contents in Spotify by ID.",
	Long: `You can like the contents in Spotify by ID.

You can like the content(s) for searching below.
  * artist
  * album
  * track`,

	// RunE : like command
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// init : executed before like command executed
func init() {
	rootCmd.AddCommand(likeCmd)

	// validate flag ''
	searchCmd.Flags().StringVarP(nil, "", "", "", "")
	if err := likeCmd.MarkFlagRequired(""); err != nil {
		return
	}

	// validate flag ''
	searchCmd.Flags().StringVarP(nil, "", "", "", "")
	if err := likeCmd.MarkFlagRequired(""); err != nil {
		return
	}
}
