package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/fevse/todo_list/internal/app"
	"github.com/fevse/todo_list/internal/storage"
)

type Server struct {
	Server *http.Server
	App    *app.App
}

func NewServer(app *app.App, host, port string) *Server {
	return &Server{
		Server: &http.Server{
			Addr:              net.JoinHostPort(host, port),
			ReadHeaderTimeout: 2 * time.Second,
		},
		App: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	m := http.NewServeMux()
	m.Handle("GET /", s.index())
	m.Handle("GET /tasks", s.ShowList())
	m.Handle("GET /tasks/{id}", s.ShowTask())
	m.Handle("POST /tasks", s.CreateTask())
	m.Handle("DELETE /tasks/{id}", s.DeleteTask())
	s.Server.Handler = m
	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *Server) index() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Hello, user!\n"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		var task storage.Task
		json.Unmarshal(body, &task)
		if err = s.App.Storage.CreateTask(task.Title, task.Status); err != nil {
			fmt.Println(err)
		}
		_, err = w.Write([]byte("Task " + task.Title + " created successfully"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			fmt.Println(err)
		}
		err = s.App.Storage.DeleteTask(id)
		if err != nil {
			fmt.Println(err)
		}

		_, err = w.Write([]byte("Task " + strconv.Itoa(id) + " successfully deleted"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) ShowList() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		list, err := s.App.Storage.ShowList()
		if err != nil {
			fmt.Println(err)
		}
		_, err = w.Write([]byte("Tasks:\n"))
		if err != nil {
			fmt.Println(err)
		}
		for _, t := range list {
			task := fmt.Sprintf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
			_, err = w.Write([]byte(task))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (s *Server) ShowTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			fmt.Println(err)
		}
		t, err := s.App.Storage.ShowTask(id)
		if err != nil {
			fmt.Println(err)
		}
		task := fmt.Sprintf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
		_, err = w.Write([]byte(task))
		if err != nil {
			fmt.Println(err)
		}
	}
}
