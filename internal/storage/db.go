package storage

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
	"gorm.io/gorm"
)

type DbStorage struct {
	connector db.Connector
}

func NewDbStorage(connector db.Connector) *DbStorage {
	return &DbStorage{
		connector: connector,
	}
}

func (s *DbStorage) Get(id uint64, contexts ...context.Context) (*Record, error) {
	session := s.getSession(contexts)
	var r Record
	err := session.First(&r, id).Error
	return &r, err
}

func (s *DbStorage) Insert(entity *Record, contexts ...context.Context) (*Record, error) {
	session := s.getSession(contexts)
	err := session.Create(entity).Error
	return entity, err
}

func (s *DbStorage) Update(entity *Record, contexts ...context.Context) (*Record, error) {
	return nil, errors.NonImplemented
}

func (s *DbStorage) Delete(id uint64, contexts ...context.Context) error {
	return errors.NonImplemented
}

func (s *DbStorage) Find(search string, contexts ...context.Context) (*Record, error) {
	return nil, errors.NonImplemented
}

func (s *DbStorage) Close() error {
	return s.connector.Close()
}

func (s *DbStorage) getSession(c []context.Context) *gorm.DB {
	var session *gorm.DB
	if len(c) > 0 {
		session = s.connector.GetConnection().Session(&gorm.Session{
			Context: c[0],
		})
	} else {
		session = s.connector.GetConnection()
	}

	return session
}
