package session

import (
	"context"

	"github.com/jmoiron/sqlx"
	handlers "github.com/mikalai-mitsin/example/internal/app/session/handlers/grpc"
	"github.com/mikalai-mitsin/example/internal/app/session/interceptors"
	"github.com/mikalai-mitsin/example/internal/app/session/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/session/usecases"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db                 *sqlx.DB
	grpcServer         *grpc.Server
	logger             *log.Log
	sessionRepository  *postgres.SessionRepository
	sessionUseCase     *usecases.SessionUseCase
	sessionInterceptor *interceptors.SessionInterceptor
	sessionHandler     *handlers.SessionServiceServer
}

func NewApp(
	db *sqlx.DB,
	grpcServer *grpc.Server,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv4Generator,
) *App {
	sessionRepository := postgres.NewSessionRepository(db, logger)
	sessionUseCase := usecases.NewSessionUseCase(sessionRepository, clock, logger, uuidGenerator)
	sessionInterceptor := interceptors.NewSessionInterceptor(sessionUseCase, logger)
	sessionHandler := handlers.NewSessionServiceServer(sessionInterceptor, logger)
	return &App{
		db:                 db,
		grpcServer:         grpcServer,
		logger:             logger,
		sessionRepository:  sessionRepository,
		sessionUseCase:     sessionUseCase,
		sessionInterceptor: sessionInterceptor,
		sessionHandler:     sessionHandler,
	}
}
func (a *App) Start(ctx context.Context) error {
	a.grpcServer.AddHandler(&examplepb.SessionService_ServiceDesc, a.sessionHandler)
	return nil
}
func (a *App) Stop(ctx context.Context) error {
	return nil
}
