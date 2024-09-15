package shs

import (
	"context"
	"errors"
	"github.com/IgorViskov/go_33_shortener/internal/algo"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
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

func (s *ShortenerService) Short(context context.Context, url string) (string, error) {
	rec, err := s.repository.Insert(context, &storage.Record{
		Value: url,
		Date:  time.Now(),
	})
	if rec == nil {
		return "", err
	}
	short := algo.Encode(rec.ID)
	return short, err
}

func (s *ShortenerService) BatchShort(context context.Context, batch []models.ShortenBatchItemDto) ([]models.ShortBatchItemDto, error) {
	dtos := ex.ToMap(batch, func(v models.ShortenBatchItemDto) string { return v.OriginalURL })
	records := ex.Map(batch, func(so models.ShortenBatchItemDto) *storage.Record {
		return &storage.Record{
			Value: so.OriginalURL,
			Date:  time.Now(),
		}
	})
	entities, errs := s.repository.BatchGetOrInsert(context, records)
	var err error
	if len(errs) > 0 {
		err = errors.Join(errs...)
	}

	result := ex.Map(entities, func(r *storage.Record) models.ShortBatchItemDto {
		return models.ShortBatchItemDto{
			CorrelationID: dtos[r.Value].CorrelationID,
			ShortURL:      algo.Encode(r.ID),
		}
	})
	return result, err
}

func (s *ShortenerService) UnShort(context context.Context, token string) (string, error) {
	id := algo.Decode(token)
	val, err := s.repository.Get(context, id)
	if err != nil {
		return "", err
	}
	return val.Value, nil
}
