package app

import (
	"github.com/fevse/todo_list/internal/logger"
	"github.com/fevse/todo_list/internal/storage"
)

type App struct {
	Storage Storage
	Logger  *logger.Logger
}

type Storage interface {
	CreateTask(string, string) error
	// UpdateTask(int64, *storage.Task) error
	DeleteTask(int) error
	ShowTask(int) (storage.Task, error)
	ShowList() ([]storage.Task, error)
	Migrate() error
}

func New(s Storage, l *logger.Logger) *App {
	return &App{
		Storage: s,
		Logger:  l,
	}
}
