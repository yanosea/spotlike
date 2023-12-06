package api

import (
	"encoding/base64"
	"errors"
	"net/url"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// constants
const (
	// error_message_invalid_length_for_random_string  is the message for failure generating random string
	error_message_invalid_length_for_random_string = "Invalid length..."
)

// getPortStringFromUri  returns port string like ':xxxx' from a URI.
func getPortStringFromUri(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if u.Port() != "" {
		return ":" + u.Port(), nil
	} else {
		return "", err
	}
}

// generateRandomString generates a random string of the specified length.
func generateRandomString(length int) (string, error) {
	if length < 0 {
		return "", errors.New(error_message_invalid_length_for_random_string)
	}

	bytes := make([]byte, length)
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
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
