package app

import "github.com/fevse/todo_list/internal/storage"

type App struct {
	Storage Storage
}

type Storage interface {
	CreateTask(string, string) error
	// UpdateTask(int64, *storage.Task) error
	DeleteTask(int64) error
	// ShowTask(int) error
	ShowList() ([]storage.Task, error)
}

func New(s Storage) *App {
	return &App{Storage: s}
}
