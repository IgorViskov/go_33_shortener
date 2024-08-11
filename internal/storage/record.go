package storage

import "time"

type Record struct {
	Id    uint64
	Value string
	Date  time.Time
}
