package storage

import "context"

type Repository[tid comparable, tentity any] interface {
	Get(id tid, context ...context.Context) (*tentity, error)
	Insert(entity *tentity, context ...context.Context) (*tentity, error)
	Update(entity *tentity, context ...context.Context) (*tentity, error)
	Delete(id tid, context ...context.Context) error
	Find(search string, context ...context.Context) (*tentity, error)
	Close() error
}
