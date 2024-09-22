package storage

import "context"

type Repository[tid comparable, tentity any] interface {
	Get(context context.Context, id tid) (*tentity, error)
	Insert(context context.Context, entity *tentity) (*tentity, error)
	Update(context context.Context, entity *tentity) (*tentity, error)
	Delete(context context.Context, id tid) error
	Find(context context.Context, search string) (*tentity, error)
	Close() error
}

type RecordRepository interface {
	Repository[uint64, Record]
	BatchGetOrInsert(context context.Context, entities []*Record) ([]*Record, []error)
	BulkDelete(context context.Context, records []*Record) error
}

type UserRepository interface {
	Repository[uint64, User]
	GetFull(context context.Context, id uint64) (*User, error)
}
