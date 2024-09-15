package apperrors

import (
	"errors"
)

var (
	ErrNonImplemented      = errors.New("not implemented")
	ErrComparatorNotFound  = errors.New("comparator not found")
	ErrInsertConflict      = errors.New("insert conflict")
	ErrRedirectUrlNotFound = errors.New("redirect url not found")
	ErrRecordNotFound      = errors.New("record not found")
	ErrInvalidUrl          = errors.New("invalid url")
	ErrInvalidJson         = errors.New("invalid json")
)
