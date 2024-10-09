package authhelper

import (
	"github.com/yanosea/spotlike/app/proxy/promptui"
)

// AuthHelpable is an interface for AuthHelper.
type AuthHelpable interface {
	AskAuthInfo() (string, string, string, string, error)
}

// AuthHelper is a struct that implements AuthHelpable.
type AuthHelper struct {
	PromptuiProxy promptuiproxy.Promptui
}

// New is a constructor of AuthHelper.
func New(promptuiProxy promptuiproxy.Promptui) *AuthHelper {
	return &AuthHelper{
		PromptuiProxy: promptuiProxy,
	}
}

func (a *AuthHelper) AskAuthInfo() (string, string, string, string, error) {
	var (
		err          error
		clientId     string
		clientSecret string
		redirectUri  string
		refreshToken string
	)

	// ask client id
	spotifyIdPrompt := a.PromptuiProxy.NewPrompt()
	spotifyIdPrompt.SetLabel(AUTH_PROMPT_SPOTIFY_ID)
	for {
		clientId, err = spotifyIdPrompt.Run()
		if err != nil {
			return "", "", "", "", err
		}
		if clientId != "" {
			break
		}
		// TODO : notify that it is required
	}
	// ask client secret
	spotifySecretPrompt := a.PromptuiProxy.NewPrompt()
	spotifySecretPrompt.SetLabel(AUTH_PROMPT_SPOTIFY_SECRET)
	for {
		clientSecret, err = spotifySecretPrompt.Run()
		if err != nil {
			return "", "", "", "", err
		}
		if clientSecret != "" {
			break
		}
		// TODO : notify that it is required
	}
	// ask redirect uri
	redirectUriPrompt := a.PromptuiProxy.NewPrompt()
	redirectUriPrompt.SetLabel(AUTH_PROMPT_SPOTIFY_REDIRECT_URI)
	for {
		redirectUri, err = redirectUriPrompt.Run()
		if err != nil {
			return "", "", "", "", err
		}
		if redirectUri != "" {
			break
		}
		// TODO : notify that it is required
	}
	// ask refresh token
	refreshTokenPrompt := a.PromptuiProxy.NewPrompt()
	refreshTokenPrompt.SetLabel(AUTH_PROMPT_SPOTIFY_REFRESH_TOKEN)
	refreshToken, err = refreshTokenPrompt.Run()
	if err != nil {
		return "", "", "", "", err
	}

	return clientId, clientSecret, redirectUri, refreshToken, nil
}
