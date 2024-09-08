package db

import "github.com/jackc/pgx/v5"

type Row interface {
	Bind(interface{}) error
}

type row struct {
	data pgx.Row
}

func (r *row) Bind(value interface{}) error {
	return r.data.Scan(value)
}
