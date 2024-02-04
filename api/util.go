package api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// constants
const (
	// error_message_invalid_uri  is the message for failure getting port from uri
	error_message_invalid_uri = "Invalid URI..."
	// error_message_invalid_length_for_random_string  is the message for failure generating random string
	error_message_invalid_length_for_random_string = "Invalid length..."
)

// getPortFromUri returns port from a URI.
func getPortFromUri(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if u.Port() == "" {
		return "", errors.New(error_message_invalid_uri)
	}

	return u.Port(), nil
}

// generateRandomString generates a random string of the specified length.
func generateRandomString(length int) (string, error) {
	if length < 0 {
		return "", errors.New(error_message_invalid_length_for_random_string)
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

// combineArtistNames concatenates the names of multiple Spotify artists into a single string.
func combineArtistNames(artists []spotify.SimpleArtist) string {
	var artistNames string
	for index, artist := range artists {
		artistNames += artist.Name
		if index+1 != len(artists) {
			artistNames += ", "
		}
	}

	return artistNames
}
