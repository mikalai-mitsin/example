package auth

import (
	"github.com/jmoiron/sqlx"
	handlers "github.com/mikalai-mitsin/example/internal/app/auth/handlers/grpc"
	"github.com/mikalai-mitsin/example/internal/app/auth/repositories/jwt"
	"github.com/mikalai-mitsin/example/internal/app/auth/services"
	"github.com/mikalai-mitsin/example/internal/app/auth/usecases"
	"github.com/mikalai-mitsin/example/internal/app/user/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/configs"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db             *sqlx.DB
	grpcServer     *grpc.Server
	logger         *log.Log
	authRepository *jwt.AuthRepository
	authService    *services.AuthService
	authUseCase    *usecases.AuthUseCase
	authHandler    *handlers.AuthServiceServer
	authMiddleware *handlers.AuthMiddleware
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
	authService := services.NewAuthService(authRepository, userRepository, logger)
	authUseCase := usecases.NewAuthUseCase(authService, clock, logger)
	authHandler := handlers.NewAuthServiceServer(authUseCase)
	authMiddleware := handlers.NewAuthMiddleware(authService, logger, config)
	return &App{
		db:             db,
		grpcServer:     grpcServer,
		logger:         logger,
		authRepository: authRepository,
		authService:    authService,
		authUseCase:    authUseCase,
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.AuthService_ServiceDesc, a.authHandler)
	grpcServer.AddInterceptor(a.authMiddleware.UnaryServerInterceptor)
	return nil
}
