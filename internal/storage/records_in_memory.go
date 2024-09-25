package storage

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/concurrent"
	"sync/atomic"
)

type InMemoryRecordStorage struct {
	current atomic.Uint64
	storage *concurrent.SyncMap[uint64, *Record]
}

func NewInMemoryRecordStorage() *InMemoryRecordStorage {
	s := &InMemoryRecordStorage{storage: concurrent.NewSyncMap[uint64, *Record]()}
	s.current.Add(1000)
	return s
}

func (i *InMemoryRecordStorage) Get(_ context.Context, id uint64) (*Record, error) {
	val, ok := i.storage.Get(id)
	if !ok {
		return nil, apperrors.ErrRedirectURLNotFound
	}
	return val, nil
}

func (i *InMemoryRecordStorage) Insert(_ context.Context, entity *Record) (*Record, error) {
	var id uint64
	exist, added := i.storage.TryAdd(entity, func() uint64 {
		id = i.current.Add(1)
		return id
	}, func(r1 *Record, r2 *Record) bool {
		return r1.Value == r2.Value
	})
	if !added {
		return exist, apperrors.ErrInsertConflict
	}
	entity.ID = id
	return entity, nil
}

func (i *InMemoryRecordStorage) BatchGetOrInsert(context context.Context, entities []*Record) ([]*Record, []error) {
	result := make([]*Record, 0, len(entities))
	err := make([]error, 0, len(entities))
	for _, e := range entities {

		added, e := i.Insert(context, e)
		if e != nil {
			err = append(err, e)
		} else {
			result = append(result, added)
		}
	}

	return result, err
}

func (i *InMemoryRecordStorage) Update(_ context.Context, entity *Record) (*Record, error) {
	i.storage.Set(entity.ID, entity)
	return entity, nil
}

func (i *InMemoryRecordStorage) Delete(_ context.Context, id uint64) error {
	r, ok := i.storage.Get(id)
	if !ok {
		return apperrors.ErrRecordNotFound
	}
	r.IsDeleted = 1
	return nil
}

func (i *InMemoryRecordStorage) BulkDelete(_ context.Context, records []*Record) error {
	for _, record := range records {
		r, _ := i.storage.Get(record.ID)
		r.IsDeleted = 1
	}

	return nil
}

func (i *InMemoryRecordStorage) Find(_ context.Context, search string) (*Record, error) {
	exist, ok := i.storage.Find(&Record{Value: search}, func(f *Record, s *Record) bool {
		return f.Value == s.Value
	})
	if !ok {
		return nil, apperrors.ErrRecordNotFound
	}
	val, _ := i.storage.Get(*exist)
	return val, nil
}

func (i *InMemoryRecordStorage) Close() error {
	return nil
}
