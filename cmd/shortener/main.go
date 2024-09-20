package main

import (
	"flag"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api"
	"github.com/IgorViskov/go_33_shortener/internal/closer"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/log"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db/migrator"
	"github.com/caarlos0/env/v11"
	"net/url"
	"os"
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
		connector := db.NewConnector(conf)
		s := selectStorage(connector, conf)
		cb.AddConfig(conf).
			UseCompression().
			Use(log.Logging()).
			AddController(app.NewShortController(conf, s)).
			AddController(app.NewUnShortController(conf, s)).
			AddController(api.NewShortenAPIController(conf, s)).
			AddController(api.NewPingAPIController(connector)).
			AddController(api.NewShortenBatchAPIController(conf, s)).
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
		CacheSize:       10,
	}

	readFlags(conf)
	if err := readEnvironments(conf); err != nil {
		log.Error(err)
	}

	return conf
}

func readFlags(conf *config.AppConfig) {

	flag.Func("a", "Адрес запуска HTTP-сервера", config.HostNameParser(conf))
	flag.Func("b", "Базовый адрес результирующего сокращённого URL", config.RedirectAddressParser(conf))
	flag.Func("f", "Путь до файла с сохраненными адресами", config.StorageFileParser(conf))
	flag.Func("d", "DSN подключения postgres", config.ConnectionStringParser(conf))
	// запускаем парсинг
	flag.Parse()
}

func readEnvironments(conf *config.AppConfig) error {
	return env.Parse(conf)
}

func selectStorage(connector db.Connector, conf *config.AppConfig) storage.Repository[uint64, storage.Record] {
	if connector.IsConnected() {
		if err := migrator.AutoMigrate(connector); err != nil {
			log.Error(err)
		}
		return storage.NewDBStorage(connector)
	}
	if conf.StorageFile != "" {
		s, err := storage.NewHybridStorage(conf)
		if err != nil {
			panic(err)
		}
		return s
	}

	return storage.NewInMemoryStorage()
}
