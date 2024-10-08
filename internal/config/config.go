package config

import (
	"net/url"
	"os"
)

type AppConfig struct {
	HostName         string  `validate:"hostname_port" env:"SERVER_ADDRESS"`
	RedirectAddress  url.URL `env:"BASE_URL"`
	StorageFile      string  `env:"FILE_STORAGE_PATH"`
	ConnectionString string  `env:"DATABASE_DSN"`
	CacheSize        int     `env:"CACHE_SIZE"`
	SecretKey        string  `env:"SECRET_KEY"`
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
		if err := tryCreateFile(flagValue); err != nil {
			return err
		}
		conf.StorageFile = flagValue
		return nil
	}
}

func ConnectionStringParser(conf *AppConfig) func(flagValue string) error {
	return func(flagValue string) error {
		conf.ConnectionString = flagValue
		return nil
	}
}

func SecretKeyParser(conf *AppConfig) func(flagValue string) error {
	return func(flagValue string) error {
		conf.SecretKey = flagValue
		return nil
	}
}

// Функция проверяет можем ли мы этом каталоге создать файл для чтения записи, если его не существует
// (есть лу у нас права, существует ли устройство и тд..)
func tryCreateFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
