package storage

import (
	"github.com/IgorViskov/go_33_shortener/internal/concurrent"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
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

func (i *InMemoryStorage) Get(id uint64) (*Record, error) {
	val, ok := i.storage.Get(id)
	if !ok {
		return nil, errors.RiseError("Redirect URL not found")
	}
	return val, nil
}

func (i *InMemoryStorage) Insert(entity *Record) (*Record, error) {
	id := i.current.Add(1)
	entity.ID = id
	i.storage.Set(id, entity)
	return entity, nil
}

func (i *InMemoryStorage) Update(entity *Record) (*Record, error) {
	i.storage.Set(entity.ID, entity)
	return entity, nil
}

func (i *InMemoryStorage) Delete(id uint64) error {
	i.storage.Remove(id)
	return nil
}

func (i *InMemoryStorage) Find(search string) (*Record, error) {
	exist, ok := i.storage.Find(&Record{Value: search}, func(f *Record, s *Record) bool {
		return f.Value == s.Value
	})
	if !ok {
		return nil, errors.RiseError("Record not found")
	}
	val, _ := i.storage.Get(*exist)
	return val, nil
}
