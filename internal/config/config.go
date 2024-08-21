package config

import (
	"net/url"
	"os"
)

type AppConfig struct {
	HostName        string  `validate:"hostname_port" env:"SERVER_ADDRESS"`
	RedirectAddress url.URL `env:"BASE_URL"`
	StorageFile     string  `json:"FILE_STORAGE_PATH"`
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
		conf.RedirectAddress = *u
		return nil
	}
}

func StorageFileParser(conf *AppConfig) func(flagValue string) error {
	return func(flagValue string) error {
		if err := exists(flagValue); err != nil {
			return err
		}
		conf.StorageFile = flagValue
		return nil
	}
}

func exists(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
