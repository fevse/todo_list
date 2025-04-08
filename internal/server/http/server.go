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
	s.App.Logger.Logger.Info("server started %v" + s.Server.Addr)
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.App.Logger.Logger.Info("server stopped")
	return s.Server.Shutdown(ctx)
}

func (s *Server) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, user!\n"))
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		s.App.Logger.Logger.Info("handler index, method " + r.Method)
	}
}

func (s *Server) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		var task storage.Task
		if err = json.Unmarshal(body, &task); err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		if err = s.App.Storage.CreateTask(task.Title, task.Status); err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		_, err = w.Write([]byte("Task " + task.Title + " created successfully"))
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		s.App.Logger.Logger.Info("Task " + task.Title + " created")
	}
}

func (s *Server) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		err = s.App.Storage.DeleteTask(id)
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}

		_, err = w.Write([]byte("Task " + strconv.Itoa(id) + " successfully deleted"))
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		s.App.Logger.Logger.Info("Task #" + strconv.Itoa(id) + " deleted")
	}
}

func (s *Server) ShowList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := make(map[string]string)
		var limit, offset int
		for k, v := range r.URL.Query() {
			if k != "limit" && k != "offset" {
				filter[k] = v[0]
			}
		}
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
			limit = -1
		}
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
			offset = 0
		}

		list, err := s.App.Storage.ShowList(filter, limit, offset)
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		_, err = w.Write([]byte("Tasks:\n"))
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		for _, t := range list {
			task := fmt.Sprintf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
			_, err = w.Write([]byte(task))
			if err != nil {
				s.App.Logger.Logger.Error(err.Error())
			}
		}
		s.App.Logger.Logger.Info("handler ShowList, method " + r.Method)
	}
}

func (s *Server) ShowTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		t, err := s.App.Storage.ShowTask(id)
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		task := fmt.Sprintf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
		_, err = w.Write([]byte(task))
		if err != nil {
			s.App.Logger.Logger.Error(err.Error())
		}
		s.App.Logger.Logger.Info("handler ShowTask, method " + r.Method)
	}
}
