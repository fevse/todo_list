package storage

import "time"

type Task struct {
	ID      int64
	Title   string
	Status  string
	Created time.Time
}
