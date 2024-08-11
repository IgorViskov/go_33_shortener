package main

import (
	"github.com/IgorViskov/go_33_shortener/cmd/shortener/api"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"net/http"
	"net/url"
)

func main() {
	mux := http.NewServeMux()
	conf := getConfig()

	builder := appInstall(conf)

	builder.Build(mux)

	err := http.ListenAndServe(conf.BaseAddress, mux)
	if err != nil {
		panic(err)
	}
}

func appInstall(conf *config.AppConfig) app.ControllerBuilder {
	builder := app.ControllerBuilder{}
	builder.AddController(api.NewMainController(conf))

	return builder
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
