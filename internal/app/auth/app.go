package auth

import (
	"github.com/jmoiron/sqlx"
	grpcHandlers "github.com/mikalai-mitsin/example/internal/app/auth/handlers/grpc"
	httpHandlers "github.com/mikalai-mitsin/example/internal/app/auth/handlers/http"
	"github.com/mikalai-mitsin/example/internal/app/auth/repositories/jwt"
	"github.com/mikalai-mitsin/example/internal/app/auth/services"
	"github.com/mikalai-mitsin/example/internal/app/auth/usecases"
	"github.com/mikalai-mitsin/example/internal/app/user/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/configs"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db                 *sqlx.DB
	grpcServer         *grpc.Server
	logger             *log.Log
	authRepository     *jwt.AuthRepository
	authService        *services.AuthService
	authUseCase        *usecases.AuthUseCase
	grpcAuthHandler    *grpcHandlers.AuthServiceServer
	httpAuthHandler    *httpHandlers.AuthHandler
	grpcAuthMiddleware *grpcHandlers.AuthMiddleware
}

func NewApp(db *sqlx.DB, config *configs.Config, logger *log.Log, clock *clock.Clock) *App {
	userRepository := postgres.NewUserRepository(db, logger)
	authRepository := jwt.NewAuthRepository(config, clock, logger)
	authService := services.NewAuthService(authRepository, userRepository, logger)
	authUseCase := usecases.NewAuthUseCase(authService, clock, logger)
	grpcAuthHandler := grpcHandlers.NewAuthServiceServer(authUseCase, logger)
	httpAuthHandler := httpHandlers.NewAuthHandler(authUseCase, logger)
	grpcAuthMiddleware := grpcHandlers.NewAuthMiddleware(authService, logger, config)
	return &App{
		db:                 db,
		logger:             logger,
		authRepository:     authRepository,
		authService:        authService,
		authUseCase:        authUseCase,
		grpcAuthHandler:    grpcAuthHandler,
		grpcAuthMiddleware: grpcAuthMiddleware,
		httpAuthHandler:    httpAuthHandler,
	}
}
func (a *App) RegisterHTTP(httpServer *http.Server) error {
	httpServer.Mount("/api/v1/auth/", a.httpAuthHandler.ChiRouter())
	return nil
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.AuthService_ServiceDesc, a.grpcAuthHandler)
	grpcServer.AddInterceptor(a.grpcAuthMiddleware.UnaryServerInterceptor)
	return nil
}
