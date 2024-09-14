package api

import (
	"encoding/json"
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
	"github.com/IgorViskov/go_33_shortener/internal/log"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
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
			return errors.RiseError("Invalid json")
		}

		shorted, err := c.service.BatchShort(dtos)

		if len(shorted) == 0 && err != nil {
			return err
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

func NewShortenBatchAPIController(config *config.AppConfig, r storage.Repository[uint64, storage.Record]) app.Controller {
	return &shortenBatchAPIController{
		path:    "/api/shorten/batch",
		service: shs.NewShortenerService(r),
		config:  config,
	}
}
