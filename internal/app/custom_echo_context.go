package app

import (
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/labstack/echo/v4"
)

type RoteContext struct {
	echo.Context
	User *storage.User
}

func GetUser(c echo.Context) *storage.User {
	return c.(*RoteContext).User
}
