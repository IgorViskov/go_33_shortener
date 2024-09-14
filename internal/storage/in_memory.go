package storage

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/concurrent"
	"sync/atomic"
)

type InMemoryStorage struct {
	current atomic.Uint64
	storage *concurrent.SyncMap[uint64, *Record]
}

func NewInMemoryStorage() *InMemoryStorage {
	s := &InMemoryStorage{storage: concurrent.NewSyncMap[uint64, *Record]()}
	s.current.Add(1000)
	return s
}

func (i *InMemoryStorage) Get(id uint64, _ ...context.Context) (*Record, error) {
	val, ok := i.storage.Get(id)
	if !ok {
		return nil, apperrors.RiseError("Redirect URL not found")
	}
	return val, nil
}

func (i *InMemoryStorage) Insert(entity *Record, _ ...context.Context) (*Record, error) {
	hashed(entity)
	var id uint64
	exist, added := i.storage.TryAdd(entity, func() uint64 {
		id = i.current.Add(1)
		return id
	}, func(r1 *Record, r2 *Record) bool {
		return r1.Hash == r2.Hash
	})
	if !added {
		return exist, apperrors.ErrInsertConflict
	}
	entity.ID = id
	return entity, nil
}

func (i *InMemoryStorage) BatchGetOrInsert(entities []*Record, contexts ...context.Context) ([]*Record, []error) {
	result := make([]*Record, 0, len(entities))
	err := make([]error, 0, len(entities))
	for _, e := range entities {

		added, e := i.Insert(e, contexts...)
		if e != nil {
			err = append(err, e)
		} else {
			result = append(result, added)
		}
	}

	return result, err
}

func (i *InMemoryStorage) Update(entity *Record, _ ...context.Context) (*Record, error) {
	i.storage.Set(entity.ID, entity)
	return entity, nil
}

func (i *InMemoryStorage) Delete(id uint64, _ ...context.Context) error {
	i.storage.Remove(id)
	return nil
}

func (i *InMemoryStorage) Find(search string, _ ...context.Context) (*Record, error) {
	exist, ok := i.storage.Find(&Record{Value: search}, func(f *Record, s *Record) bool {
		return f.Value == s.Value
	})
	if !ok {
		return nil, apperrors.RiseError("Record not found")
	}
	val, _ := i.storage.Get(*exist)
	return val, nil
}

func (i *InMemoryStorage) Close() error {
	return nil
}
