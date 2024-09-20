package app

import (
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/labstack/echo/v4"
	"net/http"
)

type unShortController struct {
	path    string
	service *shs.ShortenerService
	config  *config.AppConfig
}

func (c unShortController) Get() func(context echo.Context) error {
	return func(context echo.Context) error {
		left := len(c.config.RedirectAddress.Path) + 1
		short := context.Request().URL.Path[left:]
		moved, err := c.service.UnShort(context.Request().Context(), short)
		if err != nil {
			return err
		}
		return context.Redirect(http.StatusTemporaryRedirect, moved)
	}
}

func (c unShortController) Post() func(context echo.Context) error {
	return nil
}

func (c unShortController) GetPath() string {
	return c.path
}

func NewUnShortController(config *config.AppConfig, service *shs.ShortenerService) Controller {

	return &unShortController{
		path:    config.RedirectAddress.Path + "/*",
		service: service,
		config:  config,
	}
}
