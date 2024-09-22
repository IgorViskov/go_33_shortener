package app

import (
	"github.com/labstack/echo/v4"
)

type Controller interface {
	Get() func(c echo.Context) error
	Post() func(c echo.Context) error
	Delete() func(c echo.Context) error
	GetPath() string
}

func ErrorResult(status int, err ...error) *echo.HTTPError {
	if len(err) > 0 {
		return echo.NewHTTPError(status, err[0].Error())
	} else {
		return echo.NewHTTPError(status)
	}
}
