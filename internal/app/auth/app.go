package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
	handlers "github.com/mikalai-mitsin/example/internal/app/auth/handlers/grpc"
	"github.com/mikalai-mitsin/example/internal/app/auth/interceptors"
	"github.com/mikalai-mitsin/example/internal/app/auth/repositories/jwt"
	"github.com/mikalai-mitsin/example/internal/app/auth/usecases"
	"github.com/mikalai-mitsin/example/internal/app/user/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/configs"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db              *sqlx.DB
	grpcServer      *grpc.Server
	logger          *log.Log
	authRepository  *jwt.AuthRepository
	authUseCase     *usecases.AuthUseCase
	authInterceptor *interceptors.AuthInterceptor
	authHandler     *handlers.AuthServiceServer
}

func NewApp(
	db *sqlx.DB,
	config *configs.Config,
	grpcServer *grpc.Server,
	logger *log.Log,
	clock *clock.Clock,
) *App {
	userRepository := postgres.NewUserRepository(db, logger)
	authRepository := jwt.NewAuthRepository(config, clock, logger)
	authUseCase := usecases.NewAuthUseCase(authRepository, userRepository, logger)
	authInterceptor := interceptors.NewAuthInterceptor(authUseCase, clock, logger)
	authHandler := handlers.NewAuthServiceServer(authInterceptor)
	return &App{
		db:              db,
		grpcServer:      grpcServer,
		logger:          logger,
		authRepository:  authRepository,
		authUseCase:     authUseCase,
		authInterceptor: authInterceptor,
		authHandler:     authHandler,
	}
}
func (a *App) Start(ctx context.Context) error {
	a.grpcServer.AddHandler(&examplepb.AuthService_ServiceDesc, a.authHandler)
	return nil
}
func (a *App) Stop(ctx context.Context) error {
	return nil
}
