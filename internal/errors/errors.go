package errors

import (
	"fmt"
	"time"
)

type AppError struct {
	DateTime time.Time `json:"-"`
	Message  string    `json:"Message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("at %v, %s",
		e.DateTime, e.Message)
}

func RiseError(message string) error {
	return &AppError{
		DateTime: time.Now(),
		Message:  message,
	}
}
