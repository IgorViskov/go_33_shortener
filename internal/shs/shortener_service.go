package shs

import (
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
	short := Encode(exist.ID)
	return short, nil
}

func (s *ShortenerService) UnShort(token string) (string, error) {
	id := Decode(token)
	val, err := s.repository.Get(id)
	if err != nil {
		return "", err
	}
	return val.Value, nil
}
