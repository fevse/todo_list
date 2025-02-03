package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fevse/todo_list/internal/app"
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
	mux := http.NewServeMux()

	mux.Handle("GET /", s.index())
	// mux.Handle("GET /list/", s.ShowList())
	fmt.Println("server is running " + s.Server.Addr)
	if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *Server) index() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("*** TO-DO List ***\n"))
		if err != nil {
			fmt.Println(err)
		}
	}
}
