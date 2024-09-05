package api

import (
	"github.com/IgorViskov/go_33_shortener/internal/app"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
)

type pingAPIController struct {
	path   string
	config *config.AppConfig
}

func (c pingAPIController) Get() func(echoContext echo.Context) error {
	return func(echoContext echo.Context) error {
		conn, err := pgx.Connect(context.Background(), c.config.ConnectionString)
		if err != nil {
			return echoContext.NoContent(http.StatusInternalServerError)
		}
		return conn.Close(context.Background())
	}
}

func (c pingAPIController) Post() func(context echo.Context) error {
	return nil
}

func (c pingAPIController) GetPath() string {
	return c.path
}

func NewPingAPIController(config *config.AppConfig) app.Controller {
	return &pingAPIController{
		path:   "/ping",
		config: config,
	}
}
