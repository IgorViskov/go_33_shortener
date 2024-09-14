package shs

import (
	"github.com/IgorViskov/go_33_shortener/internal/algo"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/errors"
	"github.com/IgorViskov/go_33_shortener/internal/ex"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"time"
)

type ShortenerService struct {
	repository storage.Repository[uint64, storage.Record]
}

func NewShortenerService(r storage.Repository[uint64, storage.Record]) *ShortenerService {
	return &ShortenerService{
		repository: r,
	}
}

func (s *ShortenerService) Short(url string) (string, error) {
	exist, err := s.repository.Find(url)
	if err != nil {
		exist, err = s.repository.Insert(&storage.Record{
			Value: url,
			Date:  time.Now(),
		})
		if err != nil {
			return "", err
		}
	}
	short := algo.Encode(exist.ID)
	return short, nil
}

func (s *ShortenerService) BatchShort(batch []models.ShortenBatchItemDto) ([]models.ShortBatchItemDto, error) {
	dtos := ex.ToMap(batch, func(v models.ShortenBatchItemDto) string { return v.OriginalURL })
	records := ex.Map(batch, func(so models.ShortenBatchItemDto) storage.Record {
		return storage.Record{
			Value: so.OriginalURL,
			Date:  time.Now(),
		}
	})
	entities, errs := s.repository.BatchGetOrInsert(records)
	var err error = nil
	if errs != nil && len(errs) > 0 {
		err = errors.Combine("; ", errs...)
	}

	result := ex.Map(entities, func(r storage.Record) models.ShortBatchItemDto {
		return models.ShortBatchItemDto{
			CorrelationID: dtos[r.Value].CorrelationID,
			ShortURL:      algo.Encode(r.ID),
		}
	})
	return result, err
}

func (s *ShortenerService) UnShort(token string) (string, error) {
	id := algo.Decode(token)
	val, err := s.repository.Get(id)
	if err != nil {
		return "", err
	}
	return val.Value, nil
}
