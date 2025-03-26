package storage

import (
	"fmt"
	"time"

	"github.com/fevse/todo_list/internal/config"
	_ "github.com/jackc/pgx/stdlib" // driver
	"github.com/jmoiron/sqlx"
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
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migration, set dialect error: %w", err)
	}
	if err := goose.Up(s.db.DB, s.conf.DB.Dir); err != nil {
		return fmt.Errorf("migration up error: %w", err)
	}
	return nil
}

func (s *Storage) Connect() (err error) {
	dsn := s.conf.DBConnectionString()

	s.db, err = sqlx.Connect("pgx", dsn)
	if err != nil {
		return fmt.Errorf("connection db error: %w", err)
	}
	return nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return fmt.Errorf("db is not open")
}

func (s *Storage) ShowList() ([]Task, error) {
	data := make([]Task, 0)
	err := s.db.Select(&data, `SELECT * FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("cannot select tasks from db: %w", err)
	}
	return data, nil
}

func (s *Storage) ShowTask(id int) (Task, error) {
	data := make([]Task, 0)
	err := s.db.Select(&data, `SELECT * FROM tasks WHERE id=$1`, id)
	if err != nil {
		return Task{}, fmt.Errorf("cannot select tasks from db: %w", err)
	}
	if len(data) != 1 {
		return Task{}, fmt.Errorf("error select")
	}
	return data[0], nil
}

func (s *Storage) CreateTask(title, status string) error {
	_, err := s.db.Exec(`INSERT INTO tasks(title, status, created)
		VALUES($1, $2, $3);`,
		title, status, time.Now())
	if err != nil {
		return fmt.Errorf("cannot create task: %w", err)
	}
	return nil
}

func (s *Storage) DeleteTask(id int) error {
	_, err := s.db.Exec(`DELETE FROM tasks WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("cannot delete task %d: %w", id, err)
	}
	return nil
}
