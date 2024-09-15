package storage

import "context"

type Repository[tid comparable, tentity any] interface {
	Get(context context.Context, id tid) (*tentity, error)
	Insert(context context.Context, entity *tentity) (*tentity, error)
	Update(context context.Context, entity *tentity) (*tentity, error)
	BatchGetOrInsert(context context.Context, entities []*tentity) ([]*tentity, []error)
	Delete(context context.Context, id tid) error
	Find(context context.Context, search string) (*tentity, error)
	Close() error
}
