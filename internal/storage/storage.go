package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sqlx.DB}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(_ context.Context) (err error ){
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

func (s *Storage) CreateTask(task *Task) error {
	err := s.Connect(context.Background())
	if err != nil {
		return fmt.Errorf("cannot create task: %w", err)
	}
	defer s.Close(context.Background())
	_, err = s.db.Exec(`INSERT INTO tasks(id, title, status, created)
		VALUES($1, $2, $3, $4);`,
		task.ID, task.Title, task.Status, task.Created)
	return err
}