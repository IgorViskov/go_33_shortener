package main

import (
	"flag"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/log"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/caarlos0/env/v11"
	"net/url"
)

func main() {
	conf := getConfig()
	app.Create().Configure(configurator(conf)).Build().Start(conf)
}

func configurator(conf *config.AppConfig) app.ConfigureFunc {
	return func(cb *app.ServerBuilder) {
		s := storage.NewInMemoryStorage()
		cb.Use(log.Logging())
		cb.AddController(app.NewShortController(conf, s))
		cb.AddController(app.NewUnShortController(conf, s))
		cb.AddController(api.NewShortenApiController(conf, s))
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
	}

	readFlags(conf)
	readEnvironments(conf)

	return conf
}

func readFlags(conf *config.AppConfig) {
	flag.Func("a", "Адрес запуска HTTP-сервера", config.HostNameParser(conf))
	flag.Func("b", "Базовый адрес результирующего сокращённого URL", config.RedirectAddressParser(conf))
	// запускаем парсинг
	flag.Parse()
}

func readEnvironments(conf *config.AppConfig) {
	_ = env.Parse(conf)
}
