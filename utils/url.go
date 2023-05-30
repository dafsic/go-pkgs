package utils

import (
	"fmt"
	"net/url"
)

func UrlPath(ref string) (string, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return "", fmt.Errorf("%w%s", err, LineInfo())
	}
	return u.Path, nil
}
