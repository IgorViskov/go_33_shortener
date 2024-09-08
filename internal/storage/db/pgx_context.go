package db

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/tuples"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"iter"
)

type pgxContext struct {
	pool *pgxpool.Pool
}

func NewPgxContext(connector Connector) Context {
	dbpool := connector.GetConnection()
	if !connector.IsConnected() {
		panic(connector.GetError())
	}
	return &pgxContext{
		pool: dbpool,
	}
}

func (p *pgxContext) QueryRow(c context.Context, query string, args ...any) Row {
	return &row{
		data: p.pool.QueryRow(c, query, args...),
	}
}

func (p *pgxContext) Exec(c context.Context, query string, args ...any) error {
	t, err := p.pool.Exec(c, query, args...)
	if err != nil {
		return err
	}
	if t.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (p *pgxContext) Close() error {
	p.pool.Close()
	return nil
}

func Query[T any, U interface {
	*T
	Entity
}](p Context, c context.Context, query string, args ...any) (iter.Seq[tuples.Double[U, error]], error) {
	d, ok := p.(*pgxContext)
	if !ok {
		return nil, pgx.ErrNoRows
	}

	rows, err := d.pool.Query(c, query, args...)
	if err != nil {
		return nil, err
	}
	return func(yield func(tuples.Double[U, error]) bool) {
		defer rows.Close()
		for rows.Next() {
			if _, done := c.Deadline(); done {
				return
			}
			var u U = new(T)
			params := u.Deconstruct()
			e := rows.Scan(params...)
			if e != nil {
				if !yield(tuples.Double[U, error]{Second: &e}) {
					return
				}
			}
			if !yield(tuples.Double[U, error]{First: &u}) {
				return
			}
		}
	}, nil
}
