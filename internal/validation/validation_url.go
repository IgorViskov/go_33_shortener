package validation

import (
	"net/url"
	"strings"
)

func URL(u string) (string, bool) {
	if len(strings.TrimSpace(u)) == 0 {
		return "", false
	}
	p, err := url.Parse(u)
	if err != nil || p.Scheme == "" || p.Host == "" {
		return "", false
	}
	return u, true
}
