package config

import (
	"net/url"
)

type AppConfig struct {
	HostName        string `validate:"hostname_port"`
	RedirectAddress *url.URL
}

func HostNameParser(conf *AppConfig) func(flagValue string) error {
	return func(flagValue string) error {
		conf.HostName = flagValue
		return nil
	}
}

func RedirectAddressParser(conf *AppConfig) func(flagValue string) error {
	return func(flagValue string) error {
		u, err := url.Parse(flagValue)
		if err != nil {
			return err
		}
		conf.RedirectAddress = u
		return nil
	}
}
