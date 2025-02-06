package app

import "github.com/fevse/todo_list/internal/storage"

type App struct {
	Storage Storage
}

type Storage interface {
	CreateTask(string, string) error
	// UpdateTask(int64, *storage.Task) error
	DeleteTask(int) error
	ShowTask(int) (storage.Task, error)
	ShowList() ([]storage.Task, error)
	Migrate() error
}

func New(s Storage) *App {
	return &App{Storage: s}
}
