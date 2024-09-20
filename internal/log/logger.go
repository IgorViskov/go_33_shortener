package log

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

var logInstance *wrapper

type wrapper struct {
	logger *zap.Logger
}

func Log() *wrapper {
	if logInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if logInstance == nil {
			logInstance = initLog()
		}
	}

	return logInstance
}

func Error(e error) {
	logInstance.logger.Error(e.Error())
}

func initLog() *wrapper {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic("cannot initialize zap")
	}

	return &wrapper{
		logger: logger,
	}
}

func Logging() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,

		LogMethod:       true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			Log().logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.String("latency", fmt.Sprintf("%d ms", time.Since(v.StartTime).Milliseconds())),
				zap.String("method", c.Request().Method),
			)

			Log().logger.Info("response",
				zap.Int("status", v.Status),
				zap.Int64("response_size", v.ResponseSize),
			)
			return nil
		}})
}
