package app

import (
	"errors"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/ex"
	"github.com/IgorViskov/go_33_shortener/internal/users"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
)

var compressContentTypes = []string{"plain/text", "application/json", "application/x-gzip", "text/html"}
var authCookieName = "auth"

type ServerBuilder struct {
	controllers  []Controller
	middlewares  []echo.MiddlewareFunc
	router       *echo.Echo
	app          *Instance
	conf         *config.AppConfig
	usersManager *users.Manager
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
		if deleteHandler := c.Delete(); deleteHandler != nil {
			cb.router.DELETE(c.GetPath(), deleteHandler)
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

func (cb *ServerBuilder) AddAuth(manager *users.Manager) *ServerBuilder {
	cb.usersManager = manager
	cb.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := &users.Claims{}

			cc := &RoteContext{
				Context: c,
			}

			cookie, err := c.Request().Cookie(authCookieName)
			if err != nil {

				if c.Request().Method != echo.POST {
					return next(cc)
				}

				cc.User, err = cb.usersManager.CreateUser(c.Request().Context())
				if err != nil {
					return err
				}
				claims.UserID = cc.User.ID
				cookie, err = cb.createCookie(claims)
				if err != nil {
					return err
				}
				http.SetCookie(c.Response().Writer, cookie)
			} else {
				_, err = jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
					return []byte(cb.conf.SecretKey), nil
				})
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
				cc.User, err = cb.usersManager.FindUser(c.Request().Context(), claims.UserID)
				if err != nil {
					return echo.NewHTTPError(http.StatusNotFound, err.Error())
				}
			}

			return next(cc)
		}
	})
	return cb
}

func (cb *ServerBuilder) createCookie(claims *users.Claims) (*http.Cookie, error) {
	val, err := cb.getToken(claims)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:  authCookieName,
		Value: val,
		Path:  "/",
	}, nil
}

func (cb *ServerBuilder) getToken(claims *users.Claims) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cb.conf.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
