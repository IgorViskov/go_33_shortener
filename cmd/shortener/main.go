package main

import (
	"github.com/IgorViskov/go_33_shortener/cmd/shortener/api"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"net/url"
)

func main() {
	conf := getConfig()
	app.Create().Configure(configurator(conf)).Build().Start(conf)
}

func configurator(conf *config.AppConfig) app.ConfigureFunc {
	return func(cb *app.ServerBuilder) {
		cb.AddController(api.NewMainController(conf, storage.NewInMemoryStorage()))
	}
}

func getConfig() *config.AppConfig {
	redirect := &url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}
	return &config.AppConfig{
		RedirectAddress: redirect,
		BaseAddress:     "localhost:8080",
	}
}
