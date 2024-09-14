package app

import (
	"errors"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/ex"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
)

var compressContentTypes = []string{"plain/text", "application/json", "application/x-gzip", "text/html"}

type ServerBuilder struct {
	controllers []Controller
	middlewares []echo.MiddlewareFunc
	router      *echo.Echo
	app         *AppInstance
	conf        *config.AppConfig
}

func Create() *ServerBuilder {
	return &ServerBuilder{
		controllers: make([]Controller, 0),
		middlewares: make([]echo.MiddlewareFunc, 0),
		app:         NewAppInstance(),
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

func (cb *ServerBuilder) AddCloser(c io.Closer) *ServerBuilder {
	cb.app.AddClosable(c)
	return cb
}

func (cb *ServerBuilder) AddConfig(conf *config.AppConfig) *ServerBuilder {
	cb.conf = conf
	return cb
}

func (cb *ServerBuilder) UseCompression() *ServerBuilder {
	cb.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			h := c.Request().Header.Values(echo.HeaderContentType)
			return !ex.AnyVales(&compressContentTypes, &h)
		},
	}))
	cb.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{
		Skipper: func(c echo.Context) bool {
			h := c.Request().Header.Values(echo.HeaderContentType)
			if ex.AnyVales(&compressContentTypes, &h) && (c.Request().Header.Get(echo.HeaderContentEncoding) == "gzip") {
				return false
			}
			return true
		},
	}))

	return cb
}

type Starting interface {
	Start() error
	GetEcho() *echo.Echo
	Close()
}

func (cb *ServerBuilder) Build() Starting {
	cb.router = echo.New()

	//cb.router.HTTPErrorHandler = func(err error, c echo.Context) {
	//	if echo.Response.Header().Get(echo.HeaderContentType) == echo.MIMEApplicationJavaScript {
	//		echo.Response.Write()
	//	}
	//}

	for _, m := range cb.middlewares {
		cb.router.Use(m)
	}
	for _, c := range cb.controllers {
		if get := c.Get(); get != nil {
			cb.router.GET(c.GetPath(), get)
		}
		if post := c.Post(); post != nil {
			cb.router.POST(c.GetPath(), post)
		}
	}
	cb.router.HTTPErrorHandler = customHTTPErrorHandler
	return cb
}

func (cb *ServerBuilder) Start() error {
	if cb.conf.HostName != cb.conf.RedirectAddress.Host {
		go func() {
			if err := http.ListenAndServe(cb.conf.HostName, cb.router); !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()

		if err := http.ListenAndServe(cb.conf.RedirectAddress.Host, cb.router); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
			return err
		}
	} else {
		if err := cb.router.Start(cb.conf.HostName); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func (cb *ServerBuilder) GetEcho() *echo.Echo {
	return cb.router
}

func (cb *ServerBuilder) Close() {
	err := cb.app.Close()
	log.Printf("Close app by signal")
	if err != nil {
		log.Fatal(err)
	}
}
