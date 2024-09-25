package app

import (
	"errors"
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
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
			return ErrorResult(http.StatusBadRequest, apperrors.ErrInvalidURL)
		}
		shorted, err := c.service.Short(context.Request().Context(), u, GetUser(context))

		status := http.StatusCreated
		if err != nil {
			if errors.Is(err, apperrors.ErrInsertConflict) {
				status = http.StatusConflict
			} else {
				return ErrorResult(http.StatusInternalServerError, err)
			}
		}

		redirect := c.config.RedirectAddress
		redirect.Path = fmt.Sprintf("%s/%s", redirect.Path, shorted)

		return context.String(status, redirect.String())
	}
}

func (c shortController) Delete() func(c echo.Context) error { return nil }

func (c shortController) GetPath() string {
	return c.path
}

func NewShortController(config *config.AppConfig, service *shs.ShortenerService) Controller {
	return &shortController{
		path:    "/",
		service: service,
		config:  config,
	}
}
