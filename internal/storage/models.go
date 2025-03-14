package storage

import "time"

type Task struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}
