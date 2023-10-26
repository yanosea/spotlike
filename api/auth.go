/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package api

import (
	"context"

	"golang.org/x/oauth2/clientcredentials"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
	// https://github.com/zmb3/spotify/v2/auth
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func GetClient(clientID string, clientSecret string) (*spotify.Client, error) {
	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	if token, err := config.Token(ctx); err != nil {
		return nil, err
	} else {
		return spotify.New(spotifyauth.New().Client(ctx, token)), nil
	}
}
