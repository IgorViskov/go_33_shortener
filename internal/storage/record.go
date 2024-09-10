package storage

import (
	"strconv"
	"time"
)

type Record struct {
	ID    uint64    `gorm:"column:ID;primary_key;auto_increment"`
	Value string    `gorm:"column:Value"`
	Date  time.Time `gorm:"column:Date"`
	Hash  uint32    `gorm:"column:Hash;index"`
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

func (r *Record) Deconstruct() []interface{} {
	return []interface{}{&r.ID, &r.Value, &r.Date}
}

// TableName Имя таблицы для GORM
func (Record) TableName() string {
	return "ShortenUrls"
}
