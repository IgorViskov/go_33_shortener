package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/IgorViskov/go_33_shortener/internal/algo"
	"github.com/IgorViskov/go_33_shortener/internal/apperrors"
	"github.com/IgorViskov/go_33_shortener/internal/concurrent"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"os"
	"sync/atomic"
)

type HybridRecordStorage struct {
	current atomic.Uint64
	storage *concurrent.SyncMap[uint64, *Record]
	file    *os.File
	writer  *bufio.Writer
}

func NewHybridRecordStorage(config *config.AppConfig) (*HybridRecordStorage, error) {
	s := &HybridRecordStorage{storage: concurrent.NewSyncMap[uint64, *Record]()}
	s.current.Add(1000)

	file, err := os.OpenFile(config.StorageFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	s.file = file

	err = s.load()
	if err != nil {
		return nil, err
	}

	s.writer = bufio.NewWriter(file)
	return s, nil
}

func (s *HybridRecordStorage) Get(_ context.Context, id uint64) (*Record, error) {
	val, ok := s.storage.Get(id)
	if !ok {
		return nil, apperrors.ErrRedirectURLNotFound
	}
	return val, nil
}

func (s *HybridRecordStorage) Insert(_ context.Context, entity *Record) (*Record, error) {
	var id uint64
	exist, added := s.storage.TryAdd(entity, func() uint64 {
		id = s.current.Add(1)
		return id
	}, func(r1 *Record, r2 *Record) bool {
		return r1.Value == r2.Value
	})
	if !added {
		return exist, apperrors.ErrInsertConflict
	}
	entity.ID = id
	err := s.save(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *HybridRecordStorage) BatchGetOrInsert(context context.Context, entities []*Record) ([]*Record, []error) {
	result := make([]*Record, 0, len(entities))
	err := make([]error, 0, len(entities))
	for _, e := range entities {

		added, e := s.Insert(context, e)
		if e != nil {
			err = append(err, e)
		} else {
			result = append(result, added)
		}
	}

	return result, err
}

func (s *HybridRecordStorage) Update(_ context.Context, entity *Record) (*Record, error) {
	s.storage.Set(entity.ID, entity)
	return entity, nil
}

func (s *HybridRecordStorage) Delete(_ context.Context, id uint64) error {
	r, ok := s.storage.Get(id)
	if !ok {
		return apperrors.ErrRecordNotFound
	}
	r.IsDeleted = 1
	return nil
}

func (s *HybridRecordStorage) BulkDelete(_ context.Context, records []*Record) error {
	for _, record := range records {
		r, _ := s.storage.Get(record.ID)
		r.IsDeleted = 1
	}

	return nil
}

func (s *HybridRecordStorage) Find(_ context.Context, search string) (*Record, error) {
	exist, ok := s.storage.Find(&Record{Value: search}, func(f *Record, s *Record) bool {
		return f.Value == s.Value
	})
	if !ok {
		return nil, apperrors.ErrRecordNotFound
	}

	val, _ := s.storage.Get(*exist)

	return val, nil
}

func (s *HybridRecordStorage) load() error {
	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		dto := &RecordDto{}

		err := json.Unmarshal(data, dto)
		if err != nil {
			return err
		}

		record := dto.MapToRecord()
		s.storage.Set(record.ID, record)
	}

	return nil
}

func (s *HybridRecordStorage) save(record *Record) error {
	dto := record.MapToDto(algo.Encode(record.ID))
	data, err := json.Marshal(&dto)
	if err != nil {
		return err
	}

	if _, err := s.writer.Write(data); err != nil {
		return err
	}

	if err := s.writer.WriteByte('\n'); err != nil {
		return err
	}

	return s.writer.Flush()
}

func (s *HybridRecordStorage) Close() error {
	return s.file.Close()
}
