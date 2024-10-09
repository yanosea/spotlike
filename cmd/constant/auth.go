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
	AUTH_MESSAGE_LOGIN_SPOTIFY         = "🌐 Log in to Spotify by visiting the page below in your browser."
	AUTH_MESSAGE_AUTH_SUCCESS          = "🎉 Authentication succeeded!"
	AUTH_MESSAGE_ALREADY_AUTHENTICATED = "✅ You are already authenticated and set envs!"
	AUTH_MESSAGE_SUGGEST_SET_ENV       = "💡 If you don't want spotlike to ask questions above again, execute commands below to set envs or set your profile to set those."
	AUTH_ERROR_MESSAGE_INVALID_URI     = "❌ Invalid URI... Please check your Redirect URI and try agein..."
	AUTH_ERROR_MESSAGE_FAILED          = "❌ Authentication failed..."
	AUTH_ERROR_MESSAGE_REFRESH_FAILED  = "❌ Refresh failed... Please clear your Spotify environment variables and try again..."
	AUTH_PROMPT_SPOTIFY_ID             = "🆔 Input your Spotify Client ID"
	AUTH_PROMPT_SPOTIFY_SECRET         = "🔑 Input your Spotify Client Secret"
	AUTH_PROMPT_SPOTIFY_REDIRECT_URI   = "🔗 Input your Spotify Redirect URI"
	AUTH_PROMPT_SPOTIFY_REFRESH_TOKEN  = "🔢 Input your Spotify Refresh Token if you have one (if you don't have it, leave it empty and press enter.)"
)
