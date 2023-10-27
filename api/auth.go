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

// ClientInfo : spotify client info
type ClientInfo struct {
	// ClientID : spotify client id
	ClientID string
	// ClientSecret : spotify client secret
	ClientSecret string
}

// SpotifyClient : spotify client for using spotify api
type SpotifyClient struct {
	// Client : client of spotify
	Client *spotify.Client
	// Context : context of spotify client
	Context context.Context
}

// GetClient : returns spotify client
func GetClient() (*SpotifyClient, error) {
	// get client info
	if clientInfo, err := getClientInfo(); err != nil {
		return nil, err
	} else {
		// authenticate and get spotify client
		if client, ctx, err := auth(clientInfo.ClientID, clientInfo.ClientSecret); err != nil {
			return nil, err
		} else {
			return &SpotifyClient{
				Client:  client,
				Context: ctx,
			}, nil
		}
	}
}

// gitClientInfo : returns client info from env
func getClientInfo() (*ClientInfo, error) {
	var clientInfo *ClientInfo = new(ClientInfo)
	viper.SetConfigType("env")
	viper.SetEnvPrefix("SPOTIFY")
	viper.AutomaticEnv()

	// SPOTIFY_CLIENT_ID
	if clientID := viper.GetString("CLIENT_ID"); clientID == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Client ID",
		}

		input, err := prompt.Run()
		if err != nil {
			return nil, err
		} else {
			clientInfo.ClientID = input
		}
	} else {
		clientInfo.ClientID = clientID
	}

	// SPOTIFY_CLIENT_SECRET
	if clientSecret := viper.GetString("CLIENT_SECRET"); clientSecret == "" {
		prompt := promptui.Prompt{
			Label: "Input your Spotify Client Secret",
			Mask:  '*',
		}

		input, err := prompt.Run()
		if err != nil {
			return nil, err
		} else {
			clientInfo.ClientSecret = input
		}
	} else {
		clientInfo.ClientSecret = clientSecret
	}

	return clientInfo, nil
}

// auth: authenticates client info and returns spotify client with context
func auth(clientID string, clientSecret string) (*spotify.Client, context.Context, error) {
	ctx := context.Background()

	// authenticate config
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	// authenticate
	if token, err := config.Token(ctx); err != nil {
		return nil, nil, err
	} else {
		return spotify.New(spotifyauth.New().Client(ctx, token)), ctx, nil
	}
}
