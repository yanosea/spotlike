package app

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"

	// https://github.com/zmb3/spotify/v2
	"github.com/zmb3/spotify/v2"
)

// GetPortString  returns port string like ':xxxx' from a URI.
func GetPortString(uri string) (string, error) {
	pattern := `:(\d+)`

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(uri)

	if len(match) > 1 {
		return match[0], nil
	} else {
		return "", errors.New("\n  Invalid URI")
	}
}

// GenerateRandomString generates a random string of the specified length.
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// CombineArtistNames concatenates the names of multiple Spotify artists into a single string.
func CombineArtistNames(artists []spotify.SimpleArtist) string {
	var artistNames string
	for index, artist := range artists {
		artistNames += artist.Name
		if index+1 != len(artists) {
			artistNames += ", "
		}
	}
	return artistNames
}
