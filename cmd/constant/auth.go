package constant

const (
	AUTH_USE           = "auth"
	AUTH_HELP_TEMPLATE = `🔑 Authenticate your Spotify client.

You have to authenticate your Spotify client to use spotlike at first.
spotlike will ask you to input your Client ID, Client Secret, Redirect URI, and Refresh Token.

Usage:
  spotlike auth [flags]

Flags:
  -h, --help   help for auth
`
	AUTH_MESSAGE_ALREADY_AUTHENTICATED = "✅ You are already authenticated and set envs!"
)
