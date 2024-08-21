package api

import (
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/IgorViskov/go_33_shortener/internal/validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

type shortenApiController struct {
	path    string
	service *shs.ShortenerService
	config  *config.AppConfig
}

func (c shortenApiController) Get() func(context echo.Context) error {
	return nil
}

func (c shortenApiController) Post() func(context echo.Context) error {
	return func(context echo.Context) error {
		var dto models.ShortenDto
		err := context.Bind(&dto)
		if err != nil {
			return context.String(http.StatusBadRequest, "Invalid json")
		}
		if err != nil {
			return err
		}
		u, okValidate := validation.URL(dto.Url)
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
		return context.JSON(http.StatusCreated, responseDto)
	}
}

func (c shortenApiController) GetPath() string {
	return c.path
}

func NewShortenApiController(config *config.AppConfig, r storage.Repository[uint64, storage.Record]) *shortenApiController {
	return &shortenApiController{
		path:    "/api/shorten",
		service: shs.NewShortenerService(r),
		config:  config,
	}
}
