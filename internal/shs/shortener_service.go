package shs

import (
	"context"
	"errors"
	"fmt"
	"github.com/IgorViskov/go_33_shortener/internal/algo"
	"github.com/IgorViskov/go_33_shortener/internal/app/api/models"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/IgorViskov/go_33_shortener/internal/ex"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/labstack/gommon/log"
	"time"
)

type ShortenerService struct {
	records storage.RecordRepository
	users   storage.UserRepository
	config  *config.AppConfig
}

func NewShortenerService(r storage.RecordRepository, u storage.UserRepository, config *config.AppConfig) *ShortenerService {
	return &ShortenerService{
		records: r,
		users:   u,
		config:  config,
	}
}

func (s *ShortenerService) Short(context context.Context, url string, user *storage.User) (string, error) {
	//Создаем или получаем существующую запись (если урл совпадает)
	rec, err := s.records.Insert(context, &storage.Record{
		Value: url,
		Date:  time.Now(),
	})
	if rec == nil {
		return "", err
	}

	//Превращаем ID записи всегда в один и тот же набор символов для конкретного значения
	short := algo.Encode(rec.ID)

	user.URLs = ex.Add(user.URLs, *rec)

	_, errUpdate := s.users.Update(context, user)

	if errUpdate != nil {
		log.Error(errUpdate)
	}

	return short, err
}

func (s *ShortenerService) BatchShort(context context.Context, batch []models.ShortenBatchItemDto, user *storage.User) ([]models.ShortBatchItemDto, error) {
	dtos := ex.ToMap(batch, func(v models.ShortenBatchItemDto) string { return v.OriginalURL })
	records := ex.Map(batch, func(so models.ShortenBatchItemDto) *storage.Record {
		return &storage.Record{
			Value: so.OriginalURL,
			Date:  time.Now(),
		}
	})
	entities, errs := s.records.BatchGetOrInsert(context, records)
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

	for _, r := range records {
		user.URLs = ex.Add(user.URLs, *r)
	}

	_, errUpdate := s.users.Update(context, user)

	if errUpdate != nil {
		log.Error(errUpdate)
	}

	return result, err
}

func (s *ShortenerService) UnShort(context context.Context, token string) (string, error) {
	id := algo.Decode(token)
	val, err := s.records.Get(context, id)
	if err != nil {
		return "", err
	}
	return val.Value, nil
}

func (s *ShortenerService) EncodeURL(id uint64) string {
	redirect := s.config.RedirectAddress
	redirect.Path = fmt.Sprintf("%s/%s", redirect.Path, algo.Encode(id))
	return redirect.String()
}
