package app

import (
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type ServerBuilder struct {
	controllers []Controller
	middlewares []echo.MiddlewareFunc
	router      *echo.Echo
}

func Create() *ServerBuilder {
	return &ServerBuilder{
		controllers: make([]Controller, 0),
		middlewares: make([]echo.MiddlewareFunc, 0),
	}
}

type ConfigureFunc = func(cb *ServerBuilder)

func (cb *ServerBuilder) Configure(configure ConfigureFunc) *ServerBuilder {
	configure(cb)
	return cb
}

func (cb *ServerBuilder) AddController(c Controller) *ServerBuilder {
	cb.controllers = append(cb.controllers, c)
	return cb
}

func (cb *ServerBuilder) Use(m echo.MiddlewareFunc) *ServerBuilder {
	cb.middlewares = append(cb.middlewares, m)
	return cb
}

type Starting interface {
	Start(conf *config.AppConfig)
	GetEcho() *echo.Echo
}

func (cb *ServerBuilder) Build() Starting {
	cb.router = echo.New()

	for _, m := range cb.middlewares {
		cb.router.Use(m)
	}
	for _, c := range cb.controllers {
		cb.router.GET(c.GetPath(), c.Get)
		cb.router.POST(c.GetPath(), c.Post)
	}
	cb.router.HTTPErrorHandler = customHTTPErrorHandler
	return cb
}

func (cb *ServerBuilder) Start(conf *config.AppConfig) {
	if err := cb.router.Start(conf.BaseAddress); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (cb *ServerBuilder) GetEcho() *echo.Echo {
	return cb.router
}
