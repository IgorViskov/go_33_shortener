package api

import (
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type mainController struct {
	path    string
	service *shs.ShortenerService
	config  *config.AppConfig
}

func (c mainController) Get(context echo.Context) error {
	short := context.Request().URL.Path[1:]
	moved, err := c.service.UnShort(short)
	if err != nil {
		return err
	}
	return context.Redirect(http.StatusTemporaryRedirect, moved)
}

func (c mainController) Post(context echo.Context) error {
	body, err := io.ReadAll(context.Request().Body)
	if err != nil {
		return err
	}
	u, okValidate := validateURL(string(body))
	if !okValidate {
		return errors.RiseError("Invalid URL")
	}
	shorted, err := c.service.Short(u)

	if err != nil {
		return err
	}

	redirect := c.config.RedirectAddress
	redirect.Path = shorted

	return context.String(http.StatusCreated, redirect.String())
}

func (c mainController) GetPath() string {
	return c.path
}

func NewMainController(config *config.AppConfig, r storage.Repository[uint64, storage.Record]) *mainController {
	return &mainController{
		path:    "/*",
		service: shs.NewShortenerService(r),
		config:  config,
	}
}

func validateURL(u string) (string, bool) {
	if len(strings.TrimSpace(u)) == 0 {
		return "", false
	}
	p, err := url.Parse(u)
	if err != nil || p.Scheme == "" || p.Host == "" {
		return "", false
	}
	return u, true
}
