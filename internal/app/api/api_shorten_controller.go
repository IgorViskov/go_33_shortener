package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/IgorViskov/go_33_shortener/internal/validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

type shortenAPIController struct {
	path    string
	service *shs.ShortenerService
	config  *config.AppConfig
}

func (c shortenAPIController) Get() func(context echo.Context) error {
	return nil
}

func (c shortenAPIController) Post() func(context echo.Context) error {
	return func(context echo.Context) error {
		var dto models.ShortenDto
		err := context.Bind(&dto)
		if err != nil {
			return apperrors.ErrInvalidJson
		}
		u, okValidate := validation.URL(dto.URL)
		if !okValidate {
			return apperrors.ErrInvalidUrl
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

		responseDto := new(models.ShortDto)
		responseDto.Result = redirect.String()
		context.Response().Header().Add("Content-Type", "application/json")
		context.Response().Status = status
		return json.NewEncoder(context.Response()).Encode(&responseDto)
	}
}

func (c shortenAPIController) GetPath() string {
	return c.path
}

func NewShortenAPIController(config *config.AppConfig, r storage.Repository[uint64, storage.Record]) app.Controller {
	return &shortenAPIController{
		path:    "/api/shorten",
		service: shs.NewShortenerService(r),
		config:  config,
	}
}
