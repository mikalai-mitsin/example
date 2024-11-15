package user

import (
	"github.com/jmoiron/sqlx"
	handlers "github.com/mikalai-mitsin/example/internal/app/user/handlers/grpc"
	"github.com/mikalai-mitsin/example/internal/app/user/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/user/services"
	"github.com/mikalai-mitsin/example/internal/app/user/usecases"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db             *sqlx.DB
	logger         *log.Log
	userRepository *postgres.UserRepository
	userService    *services.UserService
	userUseCase    *usecases.UserUseCase
	userHandler    *handlers.UserServiceServer
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
	userHandler := handlers.NewUserServiceServer(userUseCase, logger)
	return &App{
		db:             db,
		logger:         logger,
		userRepository: userRepository,
		userService:    userService,
		userUseCase:    userUseCase,
		userHandler:    userHandler,
	}
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.UserService_ServiceDesc, a.userRepository)
	return nil
}
