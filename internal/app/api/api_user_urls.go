package api

import (
	"context"
	"encoding/json"
	"github.com/IgorViskov/go_33_shortener/internal/algo"
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
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
			return app.ErrorResult(http.StatusUnauthorized)
		}
		if len(user.URLs) == 0 {
			return echoContext.NoContent(http.StatusNoContent)
		}
		echoContext.Response().Header().Add("Content-Type", "application/json")
		responseDto := ex.Map(user.URLs, func(record *storage.Record) models.UserURLDto {
			return models.UserURLDto{
				ShortURL:    c.service.EncodeURL(record.ID),
				OriginalURL: record.Value,
			}
		})
		return json.NewEncoder(echoContext.Response()).Encode(&responseDto)
	}
}

func (c userURLsAPIController) Post() func(context echo.Context) error {
	return nil
}

func (c userURLsAPIController) Delete() func(echoContext echo.Context) error {
	return func(echoContext echo.Context) error {
		user := app.GetUser(echoContext)
		if user == nil {
			return app.ErrorResult(http.StatusUnauthorized)
		}

		var shorts []string
		err := echoContext.Bind(&shorts)

		echoContext.Response().Header().Add("Content-Type", "application/json")

		if err != nil {
			return app.ErrorResult(http.StatusBadRequest, apperrors.ErrInvalidJSON)
		}

		del := ex.Where(user.URLs, func(record *storage.Record) bool { return ex.Include(algo.Encode(record.ID), shorts) })

		c.service.DeleteRecordsAsync(context.Background(), del)

		return echoContext.String(http.StatusAccepted, "")
	}
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
