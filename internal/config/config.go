package config

import "net/url"

type AppConfig struct {
	BaseAddress     string
	RedirectAddress *url.URL
}
