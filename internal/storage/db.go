package storage

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/concurrent"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
	"github.com/IgorViskov/go_33_shortener/internal/ex"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db/queries"
	"log"
	"sync/atomic"
)

type CachedDbStorage struct {
	db        db.Context
	current   atomic.Uint64
	cache     *concurrent.SyncMap[uint64, *Record]
	context   context.Context
	connector db.Connector
}

func NewDbStorage(connector db.Connector, config *config.AppConfig) *CachedDbStorage {
	result := &CachedDbStorage{
		db:        db.NewPgxContext(connector),
		context:   connector.GetContext(),
		cache:     concurrent.NewSyncMap[uint64, *Record](),
		connector: connector,
	}

	result.loadCache(config.CacheSize)

	return result
}

func (s *CachedDbStorage) Get(id uint64, contexts ...context.Context) (*Record, error) {
	cached, ok := s.cache.Get(id)
	if ok {
		return cached, nil
	}
	c := ex.GetContext(contexts)
	var r Record
	err := s.db.QueryRow(c, queries.SelectRecord, id).Bind(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *CachedDbStorage) Insert(entity *Record, contexts ...context.Context) (*Record, error) {
	c := ex.GetContext(contexts)
	id := s.current.Add(1)

	entity.ID = id
	err := s.db.Exec(c, queries.InsertRecord, entity.ID, entity.Value, entity.Date)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *CachedDbStorage) Update(entity *Record, contexts ...context.Context) (*Record, error) {
	return nil, errors.NonImplemented
}

func (s *CachedDbStorage) Delete(id uint64, contexts ...context.Context) error {
	return errors.NonImplemented
}

func (s *CachedDbStorage) Find(search string, _ ...context.Context) (*Record, error) {
	return nil, errors.NonImplemented
}

func (s *CachedDbStorage) loadCache(cacheSize int) {
	reader, err := db.Query[Record](s.db, s.context, queries.SelectSeveralRecentRecords, cacheSize)
	if err != nil {
		log.Fatal(err)
	}
	var lastID uint64 = 1000
	for t := range reader {
		r, err := t.Deconstruct()
		if err != nil {
			log.Fatal(*err)
		}
		s.cache.Set((**r).ID, *r)
		lastID = (**r).ID
	}
	s.current.Add(lastID)
}

func (s *CachedDbStorage) Close() error {
	return s.connector.Close()
}
