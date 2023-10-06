/*
Copyright © 2023 yanosea <myanoshi0626@gmail.com>
*/
package api

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"

	// https://github.com/zmb3/spotify
	"github.com/zmb3/spotify"
)

func GetTokens(clientID string, clientSecret string) (string, string, error) {
	ctx := context.Background()

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}

	if token, err := config.Token(ctx); err != nil {
		return "", "", err
	} else {
		return token.AccessToken, token.RefreshToken, nil
	}
}

func CheckToken(accessToken string) (bool, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/tracks/20q73dOrP7ceLGAJQVtuTq", nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return false, nil
		}
	}

	return true, nil
}

func Refresh(clientID string, clientSecret string, refreshToken string) (string, error) {
}
