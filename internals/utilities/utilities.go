package utilities

import (
	"errors"
	"regexp"
)

func NormalizeURL(rawURL string) (string, error) {
	re, err := regexp.Compile(`^[A-z]+://|/$`)
	if err != nil {
		return "", errors.New("error compiling regular expression")
	}

	return re.ReplaceAllString(rawURL, ""), nil
}
