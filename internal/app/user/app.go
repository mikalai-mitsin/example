package user

import (
	"github.com/jmoiron/sqlx"
	grpcHandlers "github.com/mikalai-mitsin/example/internal/app/user/handlers/grpc"
	httpHandlers "github.com/mikalai-mitsin/example/internal/app/user/handlers/http"
	"github.com/mikalai-mitsin/example/internal/app/user/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/user/services"
	"github.com/mikalai-mitsin/example/internal/app/user/usecases"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db              *sqlx.DB
	logger          *log.Log
	userRepository  *postgres.UserRepository
	userService     *services.UserService
	userUseCase     *usecases.UserUseCase
	httpUserHandler *httpHandlers.UserHandler
	grpcUserHandler *grpcHandlers.UserServiceServer
}

func NewApp(
	db *sqlx.DB,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv4Generator,
) *App {
	userRepository := postgres.NewUserRepository(db, logger)
	userService := services.NewUserService(userRepository, clock, logger, uuidGenerator)
	userUseCase := usecases.NewUserUseCase(userService, logger)
	httpUserHandler := httpHandlers.NewUserHandler(userUseCase, logger)
	grpcUserHandler := grpcHandlers.NewUserServiceServer(userUseCase, logger)
	return &App{
		db:              db,
		logger:          logger,
		userRepository:  userRepository,
		userService:     userService,
		userUseCase:     userUseCase,
		httpUserHandler: httpUserHandler,
		grpcUserHandler: grpcUserHandler,
	}
}
func (a *App) RegisterHTTP(httpServer *http.Server) error {
	httpServer.Mount("/api/v1/users/", a.httpUserHandler.ChiRouter())
	return nil
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.UserService_ServiceDesc, a.grpcUserHandler)
	return nil
}
