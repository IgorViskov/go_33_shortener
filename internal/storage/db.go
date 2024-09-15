package storage

import (
	"bitbucket.org/pcastools/hash"
	"context"
	"database/sql"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
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

func (s *DBStorage) Get(context context.Context, id uint64) (*Record, error) {
	session := s.getSession(context)
	var r Record
	err := session.First(&r, id).Error
	return &r, err
}

func (s *DBStorage) Insert(context context.Context, entity *Record) (*Record, error) {
	session := s.getSession(context)
	hashed(entity)
	result := session.FirstOrCreate(entity, Record{Hash: entity.Hash})
	err := result.Error
	if err != nil {
		return entity, err
	}
	if result.RowsAffected == 0 {
		err = apperrors.ErrInsertConflict
	}
	return entity, err
}

func (s *DBStorage) BatchGetOrInsert(context context.Context, entities []*Record) ([]*Record, []error) {
	session := s.getSession(context)
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

func (s *DBStorage) Update(_ context.Context, _ *Record) (*Record, error) {
	return nil, apperrors.ErrNonImplemented
}

func (s *DBStorage) Delete(_ context.Context, _ uint64) error {
	return apperrors.ErrNonImplemented
}

func (s *DBStorage) Find(context context.Context, search string) (*Record, error) {
	h := hash.String(search)
	session := s.getSession(context)
	var r Record
	err := session.Where("\"Hash\" = $1", h).First(&r).Error
	return &r, err
}

func (s *DBStorage) Close() error {
	return s.connector.Close()
}

func (s *DBStorage) getSession(c context.Context) *gorm.DB {
	return s.connector.GetConnection(c)
}
