package app

import (
	"errors"
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/IgorViskov/go_33_shortener/internal/validation"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
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
		u, okValidate := validation.URL(string(body))
		if !okValidate {
			return apperrors.ErrInvalidURL
		}
		shorted, err := c.service.Short(context.Request().Context(), u)

		status := http.StatusCreated
		if err != nil {
			if errors.Is(err, apperrors.ErrInsertConflict) {
				status = http.StatusConflict
			} else {
				return err
			}
		}

		redirect := c.config.RedirectAddress
		redirect.Path = fmt.Sprintf("%s/%s", redirect.Path, shorted)

		return context.String(status, redirect.String())
	}
}

func (c shortController) GetPath() string {
	return c.path
}

func NewShortController(config *config.AppConfig, r storage.Repository[uint64, storage.Record]) Controller {
	return &shortController{
		path:    "/",
		service: shs.NewShortenerService(r),
		config:  config,
	}
}
