package storage

import "time"

type Record struct {
	ID    uint64
	Value string
	Date  time.Time
}
