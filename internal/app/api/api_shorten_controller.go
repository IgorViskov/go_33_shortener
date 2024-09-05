package api

import (
	"encoding/json"
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
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
			return context.String(http.StatusBadRequest, "Invalid json")
		}
		u, okValidate := validation.URL(dto.URL)
		if !okValidate {
			return errors.RiseError("Invalid URL")
		}
		shorted, err := c.service.Short(u)

		if err != nil {
			return err
		}

		redirect := c.config.RedirectAddress
		redirect.Path = fmt.Sprintf("%s/%s", redirect.Path, shorted)

		responseDto := new(models.ShortDto)
		responseDto.Result = redirect.String()
		context.Response().Header().Add("Content-Type", "application/json")
		context.Response().Status = http.StatusCreated
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
