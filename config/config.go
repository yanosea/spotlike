/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package config

import (
	"github.com/yanosea/spotlike/api"

	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/spf13/viper
	"github.com/spf13/viper"
)

type Config struct {
	Client *Client
	Tokens *Tokens
}

type Client struct {
	ClientId     string
	ClientSecret string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func New() (*Config, error) {
	viper.SetConfigType("env")
	viper.SetEnvPrefix("SPOTIFY")
	viper.AutomaticEnv()

	clientID := viper.GetString("CLIENT_ID")
	if clientID == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Client ID: ",
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
			Label: "Input your Spotify Client Secret: ",
			Mask:  '*',
		}

		input, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		clientSecret = input
	}

	accessToken := viper.GetString("ACCESS_TOKEN")
	refreshToken := viper.GetString("REFRESH_TOKEN")
	if accessToken == "" || refreshToken == "" {
		if newAccessToken, newRefreshToken, err := api.GetTokens(clientID, clientSecret); err != nil {
			return nil, err
		} else {
			accessToken = newAccessToken
			refreshToken = newRefreshToken
		}
	} else {
		if result, err := api.CheckToken(accessToken); err != nil {
			return nil, err
		} else if !result {
			if newAccessToken, err := api.Refresh(clientID, clientSecret, refreshToken); err != nil {
				return nil, err
			} else {
				accessToken = newAccessToken
			}
		}
	}

	return &Config{
		Client: &Client{
			ClientId:     clientID,
			ClientSecret: clientSecret,
		},
		Tokens: &Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
