package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	config *Config
	router *chi.Mux
	server *http.Server
}

// NewServer - provide http server
//
// @title example
// @host 127.0.0.1:8000
// @BasePath /
// @version 0.0.0
// @securitydefinitions.BearerAuth BearerAuth
func NewServer(config *Config) *Server {
	router := chi.NewRouter()
	server := &http.Server{Addr: config.Address, Handler: router}
	return &Server{server: server, config: config, router: router}
}
func (s *Server) Start(_ context.Context) error {
	return s.server.ListenAndServe()
}
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
func (s *Server) Mount(path string, handler http.Handler) {
	s.router.Mount(path, handler)
}
