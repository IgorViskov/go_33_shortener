package api

import (
	"encoding/json"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/ex"
	"github.com/IgorViskov/go_33_shortener/internal/shs"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userURLsAPIController struct {
	path    string
	service *shs.ShortenerService
}

func (c userURLsAPIController) Get() func(echoContext echo.Context) error {
	return func(echoContext echo.Context) error {
		user := app.GetUser(echoContext)
		if user == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		if len(user.URLs) == 0 {
			echoContext.Response().Status = http.StatusNoContent
			return nil
		}
		echoContext.Response().Header().Add("Content-Type", "application/json")
		responseDto := ex.Map(user.URLs, func(record storage.Record) models.UserUrlDto {
			return models.UserUrlDto{
				ShortURL:    c.service.EncodeUrl(record.ID),
				OriginalURL: record.Value,
			}
		})
		return json.NewEncoder(echoContext.Response()).Encode(&responseDto)
	}
}

func (c userURLsAPIController) Post() func(context echo.Context) error {
	return nil
}

func (c userURLsAPIController) GetPath() string {
	return c.path
}

func NewUserURLsAPIController(service *shs.ShortenerService) app.Controller {
	return &userURLsAPIController{
		path:    "/api/user/urls",
		service: service,
	}
}
