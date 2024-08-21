package storage

import (
	"strconv"
	"time"
)

type Record struct {
	ID    uint64
	Value string
	Date  time.Time
}

type RecordDto struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (r *Record) MapToDto(short string) *RecordDto {
	return &RecordDto{
		UUID:        strconv.FormatUint(r.ID, 10),
		ShortURL:    short,
		OriginalURL: r.Value,
	}
}

func (r *RecordDto) MapToRecord() *Record {
	id, err := strconv.ParseUint(r.UUID, 10, 64)
	if err != nil {
		panic("Records corrupted")
	}
	return &Record{
		ID:    id,
		Value: r.OriginalURL,
		Date:  time.Now(),
	}
}
