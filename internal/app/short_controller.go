package app

import (
	"fmt"
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

type shortController struct {
	path    string
	service *shs.ShortenerService
	config  *config.AppConfig
}

func (c shortController) Get() func(context echo.Context) error {
	return nil
}

func (c shortController) Post() func(context echo.Context) error {
	return func(context echo.Context) error {
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
		redirect.Path = fmt.Sprintf("%s/%s", redirect.Path, shorted)

		return context.String(http.StatusCreated, redirect.String())
	}
}

func (c shortController) GetPath() string {
	return c.path
}

func NewShortController(config *config.AppConfig, r storage.Repository[uint64, storage.Record]) *shortController {
	return &shortController{
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
