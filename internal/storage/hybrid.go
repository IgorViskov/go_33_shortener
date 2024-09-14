package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/IgorViskov/go_33_shortener/internal/algo"
	"github.com/IgorViskov/go_33_shortener/internal/appErrors"
	"github.com/IgorViskov/go_33_shortener/internal/concurrent"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"os"
	"sync/atomic"
)

type HybridStorage struct {
	current atomic.Uint64
	storage *concurrent.SyncMap[uint64, *Record]
	file    *os.File
	writer  *bufio.Writer
}

func NewHybridStorage(config *config.AppConfig) (*HybridStorage, error) {
	s := &HybridStorage{storage: concurrent.NewSyncMap[uint64, *Record]()}
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

func (s *HybridStorage) Get(id uint64, _ ...context.Context) (*Record, error) {
	val, ok := s.storage.Get(id)
	if !ok {
		return nil, appErrors.RiseError("Redirect URL not found")
	}
	return val, nil
}

func (s *HybridStorage) Insert(entity *Record, _ ...context.Context) (*Record, error) {
	hashed(entity)
	var id uint64
	exist, added := s.storage.TryAdd(entity, func() uint64 {
		id = s.current.Add(1)
		return id
	}, func(r1 *Record, r2 *Record) bool {
		return r1.Hash == r2.Hash
	})
	if !added {
		return exist, appErrors.ErrInsertConflict
	}
	entity.ID = id
	err := s.save(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *HybridStorage) BatchGetOrInsert(entities []*Record, contexts ...context.Context) ([]*Record, []error) {
	result := make([]*Record, 0, len(entities))
	err := make([]error, 0, len(entities))
	for _, e := range entities {

		added, e := s.Insert(e, contexts...)
		if e != nil {
			err = append(err, e)
		} else {
			result = append(result, added)
		}
	}

	return result, err
}

func (s *HybridStorage) Update(entity *Record, _ ...context.Context) (*Record, error) {
	s.storage.Set(entity.ID, entity)
	return entity, nil
}

func (s *HybridStorage) Delete(id uint64, _ ...context.Context) error {
	s.storage.Remove(id)
	return nil
}

func (s *HybridStorage) Find(search string, _ ...context.Context) (*Record, error) {
	exist, ok := s.storage.Find(&Record{Value: search}, func(f *Record, s *Record) bool {
		return f.Value == s.Value
	})
	if !ok {
		return nil, appErrors.RiseError("Record not found")
	}

	val, _ := s.storage.Get(*exist)

	return val, nil
}

func (s *HybridStorage) load() error {
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

func (s *HybridStorage) save(record *Record) error {
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

func (s *HybridStorage) Close() error {
	return s.file.Close()
}
