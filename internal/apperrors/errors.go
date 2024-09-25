package apperrors

import (
	"errors"
)

var (
	ErrNonImplemented      = errors.New("not implemented")
	ErrComparatorNotFound  = errors.New("comparator not found")
	ErrInsertConflict      = errors.New("insert conflict")
	ErrRedirectURLNotFound = errors.New("redirect url not found")
	ErrRecordNotFound      = errors.New("record not found")
	ErrInvalidURL          = errors.New("invalid url")
	ErrInvalidJSON         = errors.New("invalid json")
	ErrUserNotFound        = errors.New("user not found")
	ErrRecordIsGone        = errors.New("record is gone")
)
