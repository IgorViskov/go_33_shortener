package storage

import (
	"bitbucket.org/pcastools/hash"
	"context"
	"database/sql"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
	"gorm.io/gorm"
)

type DBStorage struct {
	connector db.Connector
}

func NewDBStorage(connector db.Connector) *DBStorage {
	return &DBStorage{
		connector: connector,
	}
}

func (s *DBStorage) Get(id uint64, contexts ...context.Context) (*Record, error) {
	session := s.getSession(contexts)
	var r Record
	err := session.First(&r, id).Error
	return &r, err
}

func (s *DBStorage) Insert(entity *Record, contexts ...context.Context) (*Record, error) {
	session := s.getSession(contexts)
	hashed(entity)
	err := session.Create(entity).Error
	return entity, err
}

func (s *DBStorage) BatchGetOrInsert(entities []*Record, contexts ...context.Context) ([]*Record, []error) {
	session := s.getSession(contexts)
	session.Begin(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	err := make([]error, 0, len(entities))
	for _, entity := range entities {
		hashed(entity)
		e := session.FirstOrCreate(entity, Record{Hash: entity.Hash}).Error
		if e != nil {
			err = append(err, e)
		}
	}
	session.Commit()

	return entities, err
}

func (s *DBStorage) Update(entity *Record, contexts ...context.Context) (*Record, error) {
	return nil, errors.NonImplemented
}

func (s *DBStorage) Delete(id uint64, contexts ...context.Context) error {
	return errors.NonImplemented
}

func (s *DBStorage) Find(search string, contexts ...context.Context) (*Record, error) {
	h := hash.String(search)
	session := s.getSession(contexts)
	var r Record
	err := session.Where("\"Hash\" = $1", h).First(&r).Error
	return &r, err
}

func (s *DBStorage) Close() error {
	return s.connector.Close()
}

func (s *DBStorage) getSession(c []context.Context) *gorm.DB {
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
