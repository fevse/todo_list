package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/fevse/todo_list/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // driver
	"github.com/pressly/goose"
)

type Storage struct {
	conf config.Config
	db   *sqlx.DB
}

func New(conf config.Config) *Storage {
	return &Storage{conf: conf}
}

func (s *Storage) Migrate() error {
	s.Connect(context.Background())
	defer s.Close(context.Background())
	if err := goose.SetDialect(s.conf.DB.Name); err != nil {
		return fmt.Errorf("migration, set dialect error: %w", err)
	}
	if err := goose.Up(s.db.DB, s.conf.DB.Dir); err != nil {
		return fmt.Errorf("migration up error: %w", err)
	}
	return nil
}

func (s *Storage) Connect(_ context.Context) (err error) {
	s.db, err = sqlx.Connect(s.conf.DB.Name, "todotasks.db")
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

func (s *Storage) ShowTask(id int) (Task, error) {
	data := make([]Task, 0)
	err := s.Connect(context.Background())
	if err != nil {
		return Task{}, fmt.Errorf("cannot show tasks: %w", err)
	}
	defer s.Close(context.Background())
	err = s.db.Select(&data, `SELECT * FROM tasks WHERE id=$1`, id)
	if err != nil {
		return Task{}, fmt.Errorf("cannot select tasks from db: %w", err)
	}
	if len(data) != 1 {
		return Task{}, fmt.Errorf("error select")
	}
	return data[0], nil
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

func (s *Storage) DeleteTask(id int) error {
	err := s.Connect(context.Background())
	if err != nil {
		return fmt.Errorf("cannot delete task: %w", err)
	}
	defer s.Close(context.Background())
	_, err = s.db.Exec(`DELETE FROM tasks WHERE id=$1`, id)
	return err
}
