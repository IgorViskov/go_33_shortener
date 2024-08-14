package shs

import (
	"math"
	"strings"
)

var alphabetConfig = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Encode(n uint64) string {
	t := make([]byte, 0)

	/* Special case */
	if n == 0 {
		return string(alphabetConfig[0])
	}

	/* Map */
	for n > 0 {
		r := n % uint64(len(alphabetConfig))
		t = append(t, alphabetConfig[r])
		n /= uint64(len(alphabetConfig))
	}

	/* Reverse */
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}

	return string(t)
}

func Decode(token string) uint64 {
	r := uint64(0)
	p := float64(len(token)) - 1

	for i := 0; i < len(token); i++ {
		r += uint64(strings.Index(alphabetConfig, string(token[i]))) * uint64(math.Pow(float64(len(alphabetConfig)), p))
		p--
	}

	return r
}
