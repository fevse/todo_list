package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // driver
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(_ context.Context) (err error) {
	s.db, err = sqlx.Connect("sqlite3", "todotasks.db")
	if err != nil {
		return fmt.Errorf("connection db error: %w", err)
	}
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return fmt.Errorf("db is not open")
}

func (s *Storage) ShowList() ([]Task, error) {
	err := s.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot show tasks: %w", err)
	}
	defer s.Close(context.Background())
	data := make([]Task, 0)
	err = s.db.Select(&data, `SELECT * FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("cannot select tasks from db: %w", err)
	}
	return data, nil
}

func (s *Storage) CreateTask(title, status string) error {
	err := s.Connect(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create task: %w", err)
	}
	defer s.Close(context.Background())
	_, err = s.db.Exec(`INSERT INTO tasks(title, status, created)
		VALUES($1, $2, $3);`,
		title, status, time.Now())
	return err
}

func (s *Storage) DeleteTask(id int64) error {
	err := s.Connect(context.Background())
	if err != nil {
		return fmt.Errorf("cannot delete task: %w", err)
	}
	defer s.Close(context.Background())
	_, err = s.db.Exec(`DELETE FROM tasks WHERE id=$1`, id)
	return err
}
