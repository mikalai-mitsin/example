package grpc

import (
	"context"
	"net"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	logger        *log.Log
	server        *grpc.Server
	config        *Config
	handlers      map[*grpc.ServiceDesc]any
	unaryUseCases []grpc.UnaryServerInterceptor
}

func NewServer(logger *log.Log, config *Config) *Server {
	return &Server{
		logger:   logger,
		server:   nil,
		config:   config,
		handlers: map[*grpc.ServiceDesc]any{},
		unaryUseCases: []grpc.UnaryServerInterceptor{
			unaryErrorServerUseCase,
			otelgrpc.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(
				logger.Logger(),
				grpc_zap.WithMessageProducer(defaultMessageProducer),
			),
		},
	}
}
func (s *Server) Start(_ context.Context) error {
	s.server = grpc.NewServer(grpc.ChainUnaryInterceptor(s.unaryUseCases...))
	for sd, ss := range s.handlers {
		s.server.RegisterService(sd, ss)
	}
	reflection.Register(s.server)
	healthServer := health.NewServer()
	for service := range s.server.GetServiceInfo() {
		healthServer.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	grpc_health_v1.RegisterHealthServer(s.server, healthServer)
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}
	return s.server.Serve(listener)
}
func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}
func (s *Server) AddHandler(sd *grpc.ServiceDesc, ss any) {
	s.handlers[sd] = ss
}
func (s *Server) AddInterceptor(usecase grpc.UnaryServerInterceptor) {
	s.unaryUseCases = append(s.unaryUseCases, usecase)
}
