package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/spotlike/app/library/authorizer"
	"github.com/yanosea/spotlike/app/library/utility"
	"github.com/yanosea/spotlike/app/proxy/cobra"
	"github.com/yanosea/spotlike/app/proxy/color"
	"github.com/yanosea/spotlike/app/proxy/context"
	"github.com/yanosea/spotlike/app/proxy/http"
	"github.com/yanosea/spotlike/app/proxy/io"
	"github.com/yanosea/spotlike/app/proxy/oauth2"
	"github.com/yanosea/spotlike/app/proxy/os"
	"github.com/yanosea/spotlike/app/proxy/promptui"
	"github.com/yanosea/spotlike/app/proxy/randstr"
	"github.com/yanosea/spotlike/app/proxy/spotify"
	"github.com/yanosea/spotlike/app/proxy/spotifyauth"
	"github.com/yanosea/spotlike/app/proxy/url"
	"github.com/yanosea/spotlike/cmd/constant"
)

// authOption is the struct for auth command.
type authOption struct {
	Out           ioproxy.WriterInstanceInterface
	ErrOut        ioproxy.WriterInstanceInterface
	Args          []string
	Authorizer    authorizer.Authorizable
	PromptuiProxy promptuiproxy.Promptui
	Utility       utility.UtilityInterface
}

// NewAuthCommand creates a new auth command.
func NewAuthCommand(g *GlobalOption, promptuiProxy promptuiproxy.Promptui) *cobraproxy.CommandInstance {
	o := &authOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    g.Args,
		Utility: g.Utility,
	}
	o.Authorizer = authorizer.New(
		contextproxy.New(),
		httpproxy.New(),
		oauth2proxy.New(),
		osproxy.New(),
		randstrproxy.New(),
		spotifyproxy.New(),
		spotifyauthproxy.New(),
		urlproxy.New(),
	)
	o.PromptuiProxy = promptuiProxy

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.AUTH_USE
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.authRunE

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.AUTH_HELP_TEMPLATE)

	return cmd
}

// authRunE is the function that is called when the auth command is executed.
func (o *authOption) authRunE(_ *cobra.Command, _ []string) error {
	// check if auth info is already set
	if o.Authorizer.IsEnvsSet() {
		colorProxy := colorproxy.New()
		// if already set, print message and return
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.AUTH_MESSAGE_ALREADY_AUTHENTICATED))
		return nil
	}

	return o.auth()
}

func (o *authOption) auth() error {
	// get auth info
	clientId, clientSecret, redirectUri, refreshToken, err := o.askAuthInfo()
	if err != nil {
		return err
	}
	// set auth info
	if status := o.Authorizer.SetAuthInfo(clientId, clientSecret, redirectUri, refreshToken); status != authorizer.SetEnvSuccessfully {
		// TODO : print what was lack of
		return nil
	}
	// authenticate or refresh
	var status authorizer.AuthenticateStatus
	if refreshToken == "" {
		// if refresh token is not set, get and print auth url
		o.Utility.PrintlnWithWriter(o.Out, o.Authorizer.GetAuthUrl())
		// authenticate
		_, status, err = o.Authorizer.Authenticate()
		if err != nil {
			return err
		}
	} else {
		// if refresh token is set, refresh
		_, status = o.Authorizer.Refresh()

	}

	// TODO : print success message
	o.Utility.PrintlnWithWriter(o.Out, status)

	return nil
}

func (o *authOption) askAuthInfo() (string, string, string, string, error) {
	var (
		err          error
		clientId     string
		clientSecret string
		redirectUri  string
		refreshToken string
	)

	// ask client id
	spotifyIdPrompt := o.PromptuiProxy.NewPrompt()
	spotifyIdPrompt.SetLabel(constant.AUTH_PROMPT_SPOTIFY_ID)
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
	spotifySecretPrompt := o.PromptuiProxy.NewPrompt()
	spotifySecretPrompt.SetLabel(constant.AUTH_PROMPT_SPOTIFY_SECRET)
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
	redirectUriPrompt := o.PromptuiProxy.NewPrompt()
	redirectUriPrompt.SetLabel(constant.AUTH_PROMPT_SPOTIFY_REDIRECT_URI)
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
	refreshTokenPrompt := o.PromptuiProxy.NewPrompt()
	refreshTokenPrompt.SetLabel(constant.AUTH_PROMPT_SPOTIFY_REFRESH_TOKEN)
	refreshToken, err = refreshTokenPrompt.Run()
	if err != nil {
		return "", "", "", "", err
	}

	return clientId, clientSecret, redirectUri, refreshToken, nil
}
