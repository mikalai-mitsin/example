package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/riandyrn/otelchi"
)

type Server struct {
	config *Config
	router *chi.Mux
	server *http.Server
	logger log.Logger
}

// NewServer - provide http server
//
// @title example
// @host http://127.0.0.1:8000
// @BasePath /
// @version 0.0.0
// @securitydefinitions.BearerAuth BearerAuth
func NewServer(config *Config, logger log.Logger) *Server {
	router := chi.NewRouter()
	router.Use(otelchi.Middleware("example"))
	router.Use(loggerMiddleware(logger))
	server := &http.Server{Addr: config.Address, Handler: router}
	return &Server{server: server, config: config, router: router, logger: logger}
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
func loggerMiddleware(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.WithContext(r.Context()).
				Info("finished http request", log.String("system", "http"), log.String("http.method", r.Method), log.String("http.path", r.URL.Path), log.String("http.remote_addr", r.RemoteAddr), log.Int("http.status", ww.Status()), log.Int64("http.time_ms", time.Since(start).Milliseconds()), log.String("http.start_time", start.Format(time.RFC3339)))
		})
	}
}
