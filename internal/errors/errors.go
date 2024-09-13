package errors

import (
	"fmt"
	"strings"
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

func Combine(separator string, err ...error) error {
	var sb strings.Builder
	l := len(err)
	for i, e := range err {
		sb.WriteString(e.Error())
		if i < l-1 {
			sb.WriteString(separator)
		}
	}
	return RiseError(sb.String())
}

var (
	NonImplemented     = RiseError("Not Implemented")
	ComparatorNotFound = RiseError("Comparator not found")
)
