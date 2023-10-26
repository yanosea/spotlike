/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package client

import (
	"github.com/yanosea/spotlike/api"

	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/spf13/viper
	"github.com/spf13/viper"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

func New() (*spotify.Client, error) {
	viper.SetConfigType("env")
	viper.SetEnvPrefix("SPOTIFY")
	viper.AutomaticEnv()

	clientID := viper.GetString("CLIENT_ID")
	if clientID == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Client ID",
		}

		input, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		clientID = input
	}

	clientSecret := viper.GetString("CLIENT_SECRET")
	if clientSecret == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Client Secret",
			Mask:  '*',
		}

		input, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		clientSecret = input
	}

	if client, err := api.GetClient(clientID, clientSecret); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
