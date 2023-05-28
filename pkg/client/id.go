package client

import "regexp"

var base62Regexp = regexp.MustCompile(`^[0-9A-Za-z]+$`)

func IsSpotifyID(value string) bool {
	return base62Regexp.MatchString(value) && len(value) == 22
}
