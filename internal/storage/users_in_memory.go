package storage

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/concurrent"
	"sync/atomic"
)

type InMemoryUsersStorage struct {
	current atomic.Uint64
	storage *concurrent.SyncMap[uint64, *User]
}

func NewInMemoryUsersStorage() *InMemoryUsersStorage {
	s := &InMemoryUsersStorage{storage: concurrent.NewSyncMap[uint64, *User]()}
	return s
}

func (i *InMemoryUsersStorage) Get(_ context.Context, id uint64) (*User, error) {
	val, ok := i.storage.Get(id)
	if !ok {
		return nil, apperrors.ErrUserNotFound
	}
	return val, nil
}

func (i *InMemoryUsersStorage) GetFull(context context.Context, id uint64) (*User, error) {
	return i.Get(context, id)
}

func (i *InMemoryUsersStorage) Insert(_ context.Context, entity *User) (*User, error) {
	var id uint64
	exist, added := i.storage.TryAdd(entity, func() uint64 {
		id = i.current.Add(1)
		return id
	}, func(r1 *User, r2 *User) bool {
		return r1.ID == r2.ID
	})
	if !added {
		return exist, apperrors.ErrInsertConflict
	}
	entity.ID = id
	return entity, nil
}

func (i *InMemoryUsersStorage) Update(_ context.Context, user *User) (*User, error) {
	i.storage.Set(user.ID, user)
	return user, nil
}

func (i *InMemoryUsersStorage) Delete(_ context.Context, _ uint64) error {
	return apperrors.ErrNonImplemented
}

func (i *InMemoryUsersStorage) Find(_ context.Context, _ string) (*User, error) {
	return nil, apperrors.ErrNonImplemented
}

func (i *InMemoryUsersStorage) Close() error {
	return nil
}
