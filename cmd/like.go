package cmd

import (
	"errors"
	"fmt"

	// https://github.com/fatih/color
	"github.com/fatih/color"
	// https://github.com/spf13/cobra
	"github.com/spf13/cobra"
)

const (
	like_use   = "like"
	like_short = "🤍 Like content in Spotify by ID."

	like_long = `🤍 Like content in Spotify by ID.

You must use sub command below...

* 🖌️  artist
* 💿 album
* 🎵 track`
	like_error_message_no_sub_command = `Use sub command below...

* 🖌️  artist
* 💿 album
* 🎵 track`
	like_error_message_args_or_flag_id_required = `The arguments or the flag "id" is required...`
)

func newLikeCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   like_use,
		Short: like_short,
		Long:  like_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New(like_error_message_no_sub_command)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	cmd.AddCommand(
		newLikeArtistCommand(globalOption),
		newLikeAlbumCommand(globalOption),
		newLikeTrackCommand(globalOption),
	)

	return cmd
}
func FormatLikeResultError(error error) string {
	return color.RedString(fmt.Sprintf("Error:\n  %s", error))
}
