package db

import (
	"context"
)

type Context interface {
	QueryRow(context.Context, string, ...any) Row
	Exec(context.Context, string, ...any) error
}
