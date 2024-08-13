package app

import (
	"github.com/labstack/echo/v4"
)

type Controller interface {
	Get() func(c echo.Context) error
	Post() func(c echo.Context) error
	GetPath() string
}
