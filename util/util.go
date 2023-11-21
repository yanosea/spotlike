package util

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// GetPortStringFromUri  returns port string like ':xxxx' from a URI.
func GetPortStringFromUri(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", errors.New("\n  Invalid URI\n")
	}

	if u.Port() != "" {
		return ":" + u.Port(), nil
	} else {
		return "", errors.New("\n  Invalid URI\n")
	}
}

// GenerateRandomString generates a random string of the specified length.
func GenerateRandomString(length int) (string, error) {
	if length < 0 {
		return "", errors.New("\n  Invalid length\n")
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// CombineArtistNames concatenates the names of multiple Spotify artists into a single string.
func CombineArtistNames(artists []spotify.SimpleArtist) (string, error) {
	if len(artists) == 0 || artists == nil {
		return "", errors.New("\n  Invalid artists\n")
	}

	var artistNames string
	for index, artist := range artists {
		artistNames += artist.Name
		if index+1 != len(artists) {
			artistNames += ", "
		}
	}
	return artistNames, nil
}
