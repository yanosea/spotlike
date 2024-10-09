package spotifyauthproxy

import (
	"github.com/zmb3/spotify/v2/auth"
)

const (
	// ScopeUserFollowRead is a const for spotifyauth.ScopeUserFollowRead.
	ScopeUserFollowRead = spotifyauth.ScopeUserFollowRead
	// ScopeUserFollowModify is a const for spotifyauth.ScopeUserFollowModify.
	ScopeUserFollowModify = spotifyauth.ScopeUserFollowModify
	// ScopeUserLibraryRead is a const for spotifyauth.ScopeUserLibraryRead.
	ScopeUserLibraryRead = spotifyauth.ScopeUserLibraryRead
	// ScopeUserLibraryModify is a const for spotifyauth.ScopeUserLibraryModify.
	ScopeUserLibraryModify = spotifyauth.ScopeUserLibraryModify
)

var (
	// WithScpoes is a variable for spotifyauth.WithScopes.
	WithScopes = spotifyauth.WithScopes
	// WithState is a variable for spotifyauth.WithState.
	WithClientID = spotifyauth.WithClientID
	// WithClientSecret is a variable for spotifyauth.WithClientSecret.
	WithClientSecret = spotifyauth.WithClientSecret
	// WithRedirectURL is a variable for spotifyauth.WithRedirectURL.
	WithRedirectURL = spotifyauth.WithRedirectURL
)
