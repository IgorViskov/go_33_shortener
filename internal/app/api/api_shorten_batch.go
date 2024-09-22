package api

import (
	"encoding/json"
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/log"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/labstack/echo/v4"
	"net/http"
)

type shortenBatchAPIController struct {
	path    string
	service *shs.ShortenerService
	config  *config.AppConfig
}

func (c shortenBatchAPIController) Get() func(context echo.Context) error {
	return nil
}

func (c shortenBatchAPIController) Post() func(context echo.Context) error {
	return func(context echo.Context) error {
		var dtos []models.ShortenBatchItemDto
		err := context.Bind(&dtos)

		context.Response().Header().Add("Content-Type", "application/json")

		if err != nil {
			return app.ErrorResult(http.StatusBadRequest, apperrors.ErrInvalidJSON)
		}

		shorted, err := c.service.BatchShort(context.Request().Context(), dtos, app.GetUser(context))

		if len(shorted) == 0 && err != nil {
			return app.ErrorResult(http.StatusBadRequest, err)
		} else if err != nil {
			log.Error(err)
		}

		result := make([]models.ShortBatchItemDto, 0, len(shorted))
		for _, r := range shorted {
			redirect := c.config.RedirectAddress
			redirect.Path = fmt.Sprintf("%s/%s", redirect.Path, r.ShortURL)
			r.ShortURL = redirect.String()
			result = append(result, r)
		}

		context.Response().Header().Add("Content-Type", "application/json")
		context.Response().Status = http.StatusCreated
		return json.NewEncoder(context.Response()).Encode(&result)
	}
}

func (c shortenBatchAPIController) GetPath() string {
	return c.path
}

func (c shortenBatchAPIController) Delete() func(c echo.Context) error { return nil }

func NewShortenBatchAPIController(config *config.AppConfig, service *shs.ShortenerService) app.Controller {
	return &shortenBatchAPIController{
		path:    "/api/shorten/batch",
		service: service,
		config:  config,
	}
}
