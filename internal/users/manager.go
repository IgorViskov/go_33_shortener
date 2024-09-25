package users

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
)

type Manager struct {
	repository storage.UserRepository
}

func NewManager(repository storage.UserRepository) *Manager {
	return &Manager{repository}
}

func (m *Manager) CreateUser(context context.Context) (*storage.User, error) {
	return m.repository.Insert(context, &storage.User{})
}

func (m *Manager) FindUser(context context.Context, id uint64) (*storage.User, error) {
	return m.repository.GetFull(context, id)
}
