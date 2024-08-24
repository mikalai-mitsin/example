package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	handlers "github.com/mikalai-mitsin/example/internal/app/user/handlers/grpc"
	"github.com/mikalai-mitsin/example/internal/app/user/interceptors"
	"github.com/mikalai-mitsin/example/internal/app/user/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/user/usecases"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db              *sqlx.DB
	grpcServer      *grpc.Server
	logger          *log.Log
	userRepository  *postgres.UserRepository
	userUseCase     *usecases.UserUseCase
	userInterceptor *interceptors.UserInterceptor
	userHandler     *handlers.UserServiceServer
}

func NewApp(
	db *sqlx.DB,
	grpcServer *grpc.Server,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv4Generator,
) *App {
	userRepository := postgres.NewUserRepository(db, logger)
	userUseCase := usecases.NewUserUseCase(userRepository, clock, logger, uuidGenerator)
	userInterceptor := interceptors.NewUserInterceptor(userUseCase, logger)
	userHandler := handlers.NewUserServiceServer(userInterceptor, logger)
	return &App{
		db:              db,
		grpcServer:      grpcServer,
		logger:          logger,
		userRepository:  userRepository,
		userUseCase:     userUseCase,
		userInterceptor: userInterceptor,
		userHandler:     userHandler,
	}
}
func (a *App) Start(ctx context.Context) error {
	a.grpcServer.AddHandler(&examplepb.UserService_ServiceDesc, a.userHandler)
	return nil
}
func (a *App) Stop(ctx context.Context) error {
	return nil
}
