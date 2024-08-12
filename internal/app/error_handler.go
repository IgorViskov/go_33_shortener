package app

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusBadRequest
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.JSON(code, err)
}
