package authorizer

// SetEnvStatus is a type for set environment status
type SetEnvStatus int

const (
	// SetEnvSuccessfully is a status for set environment successfully
	SetEnvSuccessfully SetEnvStatus = iota
	// SetEnvFailedId is a status for set environment failed because lack of ID
	SetEnvFailedId
	// SetEnvFailedSecret is a status for set environment failed because lack of secret
	SetEnvFailedSecret
	// SetEnvFailedRedirectUri is a status for set environment failed because lack of redirect URI
	SetEnvFailedRedirectUri
)

// AuthenticateStatus is a type for authenticate status
type AuthenticateStatus int

const (
	// AuthenticateedSuccessfully is a status for authenticateed successfully
	AuthenticatedSuccessfully AuthenticateStatus = iota
	// AuthenticateFailedInvalidUri is a status for authenticateed failed because invalid URI
	AuthenticateFailedInvalidUri
	// AuthenticateedFailed is a status for authenticateed failed
	AuthenticateFailed
	// AuthenticateedAlready is a status for authenticateed already
	AuthenticateAlready
)

const (
	ENV_SPOTIFY_ID            = "SPOTIFY_ID"
	ENV_SPOTIFY_SECRET        = "SPOTIFY_SECRET"
	ENV_SPOTIFY_REDIRECT_URI  = "SPOTIFY_REDIRECT_URI"
	ENV_SPOTIFY_REFRESH_TOKEN = "SPOTIFY_REFRESH_TOKEN"
	TEMPLATE_SET_ENV_COMMAND  = "export %s="
)
