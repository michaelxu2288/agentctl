package api

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	handlers   *Handlers
}

func NewServer(addr string, h *Handlers) *Server {
	if addr == "" {
		addr = ":7070"
	}
	mux := http.NewServeMux()
	h.Register(mux)
	return &Server{
		httpServer: &http.Server{Addr: addr, Handler: mux},
		handlers:   h,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if ctx == nil {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ctx = c
	}
	return s.httpServer.Shutdown(ctx)
}
