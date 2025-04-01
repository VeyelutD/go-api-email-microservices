package api

import (
	"log/slog"
	"net/http"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Run(router http.Handler) error {
	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	slog.Info("Starting API server", "addr", s.addr)
	return server.ListenAndServe()
}
