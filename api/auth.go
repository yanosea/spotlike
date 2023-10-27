/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package api

import (
	"context"

	"golang.org/x/oauth2/clientcredentials"

	// https://github.com/manifoldco/promptui
	"github.com/manifoldco/promptui"
	// https://github.com/spf13/viper
	"github.com/spf13/viper"
	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
	// https://github.com/zmb3/spotify/v2/auth
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// SpotifyClient : spotify client for using spotify api
type SpotifyClient struct {
	// Client : client of spotify
	Client *spotify.Client
	// Context : context of spotify client
	Context context.Context
}

var (
	clientID     string
	clientSecret string
)

func GetClient() (*SpotifyClient, error) {
	if client, ctx, err := auth(clientID, clientSecret); err != nil {
		return nil, err
	} else {
		return &SpotifyClient{
			Client:  client,
			Context: ctx,
		}, nil
	}
}

func setClientInfo() error {
	viper.SetConfigType("env")
	viper.SetEnvPrefix("SPOTIFY")
	viper.AutomaticEnv()

	if viper.GetString("CLIENT_ID") == "" || viper.GetString("CLIENT_SECRET") == "" {
		promptID := promptui.Prompt{
			Label: "Input your Spotify Client ID",
		}

		inputID, err := promptID.Run()
		if err != nil {
			return err
		} else {
			clientID = inputID
		}

		promptSecret := promptui.Prompt{
			Label: "Input your Spotify Client Secret",
			Mask:  '*',
		}

		inputSecret, err := promptSecret.Run()
		if err != nil {
			return err
		} else {
			clientSecret = inputSecret
		}
	} else {

	}

	return nil
}

func auth(clientID string, clientSecret string) (*spotify.Client, context.Context, error) {
	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	if token, err := config.Token(ctx); err != nil {
		return nil, nil, err
	} else {
		return spotify.New(spotifyauth.New().Client(ctx, token)), ctx, nil
	}
}
