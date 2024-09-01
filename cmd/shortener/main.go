package main

import (
	"flag"
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api"
	"github.com/IgorViskov/go_33_shortener/internal/closer"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/log"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/caarlos0/env/v11"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	conf := getConfig()
	builder := app.Create().Configure(configurator(conf)).Build()
	closer.Bind(builder.Close)
	exitCode := closer.Checked(builder.Start)
	os.Exit(exitCode)
}

func configurator(conf *config.AppConfig) app.ConfigureFunc {
	return func(cb *app.ServerBuilder) {
		s, err := storage.NewHybridStorage(conf)
		if err != nil {
			panic(err)
		}
		cb.AddConfig(conf).
			UseCompression().
			Use(log.Logging()).
			AddController(app.NewShortController(conf, s)).
			AddController(app.NewUnShortController(conf, s)).
			AddController(api.NewShortenAPIController(conf, s)).
			AddCloser(s)
	}
}

func getConfig() *config.AppConfig {
	redirect := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}

	conf := &config.AppConfig{
		RedirectAddress: redirect,
		HostName:        "localhost:8080",
		StorageFile:     fmt.Sprintf("%s%c%s", getExecuteDir(), os.PathSeparator, "db.json"),
	}

	readFlags(conf)
	readEnvironments(conf)

	return conf
}

func readFlags(conf *config.AppConfig) {
	flag.Func("a", "Адрес запуска HTTP-сервера", config.HostNameParser(conf))
	flag.Func("b", "Базовый адрес результирующего сокращённого URL", config.RedirectAddressParser(conf))
	flag.Func("f", "Путь до файла с сохраненными адресами", config.StorageFileParser(conf))
	// запускаем парсинг
	flag.Parse()
}

func readEnvironments(conf *config.AppConfig) {
	_ = env.Parse(conf)
}

func getExecuteDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
