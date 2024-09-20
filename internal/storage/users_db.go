package storage

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBUsersStorage struct {
	connector db.Connector
}

func NewDBUsersStorage(connector db.Connector) *DBUsersStorage {
	return &DBUsersStorage{
		connector: connector,
	}
}

func (s *DBUsersStorage) Get(context context.Context, id uint64) (*User, error) {
	session := s.getSession(context)
	var r User
	err := session.First(&r, id).Error
	return &r, err
}

func (s *DBUsersStorage) Insert(context context.Context, entity *User) (*User, error) {
	session := s.getSession(context)
	err := session.Create(entity).Error
	return entity, err
}

func (s *DBUsersStorage) Update(context context.Context, entity *User) (*User, error) {
	session := s.getSession(context)
	err := session.Save(entity).Error
	return entity, err
}

func (s *DBUsersStorage) Delete(_ context.Context, _ uint64) error {
	return apperrors.ErrNonImplemented
}

func (s *DBUsersStorage) Find(_ context.Context, _ string) (*User, error) {
	return nil, apperrors.ErrNonImplemented
}

func (s *DBUsersStorage) GetFull(context context.Context, id uint64) (*User, error) {
	session := s.getSession(context)
	u := &User{}
	e := session.Preload(clause.Associations).Find(u, id).Error
	return u, e
}

func (s *DBUsersStorage) Close() error {
	return s.connector.Close()
}

func (s *DBUsersStorage) getSession(c context.Context) *gorm.DB {
	return s.connector.GetConnection(c)
}
