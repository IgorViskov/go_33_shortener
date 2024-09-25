package storage

import (
	"bitbucket.org/pcastools/hash"
	"context"
	"database/sql"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBRecordsStorage struct {
	connector db.Connector
}

func NewDBRecordsStorage(connector db.Connector) *DBRecordsStorage {
	return &DBRecordsStorage{
		connector: connector,
	}
}

func (s *DBRecordsStorage) Get(context context.Context, id uint64) (*Record, error) {
	session := s.getSession(context)
	var r Record
	err := session.Unscoped().First(&r, id).Error
	return &r, err
}

func (s *DBRecordsStorage) Insert(context context.Context, entity *Record) (*Record, error) {
	session := s.getSession(context)
	result := session.FirstOrCreate(entity, Record{Value: entity.Value})
	err := result.Error
	if err != nil {
		return entity, err
	}
	if result.RowsAffected == 0 {
		err = apperrors.ErrInsertConflict
	}
	return entity, err
}

func (s *DBRecordsStorage) BatchGetOrInsert(context context.Context, entities []*Record) ([]*Record, []error) {
	session := s.getSession(context)
	session.Begin(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	err := make([]error, 0, len(entities))
	for _, entity := range entities {
		e := session.FirstOrCreate(entity, Record{Value: entity.Value}).Error
		if e != nil {
			err = append(err, e)
		}
	}
	session.Commit()

	return entities, err
}

func (s *DBRecordsStorage) Update(_ context.Context, _ *Record) (*Record, error) {
	return nil, apperrors.ErrNonImplemented
}

func (s *DBRecordsStorage) Delete(context context.Context, id uint64) error {
	session := s.getSession(context)
	return session.Association(clause.Associations).Delete(&Record{ID: id})
}
func (s *DBRecordsStorage) BulkDelete(context context.Context, records []*Record) error {
	session := s.getSession(context)
	return session.Select(clause.Associations).Delete(&records).Error
}

func (s *DBRecordsStorage) Find(context context.Context, search string) (*Record, error) {
	h := hash.String(search)
	session := s.getSession(context)
	var r Record
	err := session.Where("\"Hash\" = $1", h).First(&r).Error
	return &r, err
}

func (s *DBRecordsStorage) Close() error {
	return s.connector.Close()
}

func (s *DBRecordsStorage) getSession(c context.Context) *gorm.DB {
	return s.connector.GetConnection(c)
}
