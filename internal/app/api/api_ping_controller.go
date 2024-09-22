package api

import (
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

type pingAPIController struct {
	path      string
	connector db.Connector
}

func (c pingAPIController) Get() func(echoContext echo.Context) error {
	return func(echoContext echo.Context) error {
		if c.connector.IsConnected() {
			return nil
		}
		c.connector.GetConnection(echoContext.Request().Context())
		if !c.connector.IsConnected() {
			return echoContext.NoContent(http.StatusInternalServerError)
		}
		return nil
	}
}

func (c pingAPIController) Post() func(context echo.Context) error {
	return nil
}

func (c pingAPIController) Delete() func(c echo.Context) error { return nil }

func (c pingAPIController) GetPath() string {
	return c.path
}

func NewPingAPIController(connector db.Connector) app.Controller {
	return &pingAPIController{
		path:      "/ping",
		connector: connector,
	}
}
