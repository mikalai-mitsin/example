package arch

import (
	"context"

	"github.com/jmoiron/sqlx"
	handlers "github.com/mikalai-mitsin/example/internal/app/arch/handlers/grpc"
	"github.com/mikalai-mitsin/example/internal/app/arch/interceptors"
	"github.com/mikalai-mitsin/example/internal/app/arch/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/arch/usecases"
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
	archRepository  *postgres.ArchRepository
	archUseCase     *usecases.ArchUseCase
	archInterceptor *interceptors.ArchInterceptor
	archHandler     *handlers.ArchServiceServer
}

func NewApp(
	db *sqlx.DB,
	grpcServer *grpc.Server,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv4Generator,
) *App {
	archRepository := postgres.NewArchRepository(db, logger)
	archUseCase := usecases.NewArchUseCase(archRepository, clock, logger, uuidGenerator)
	archInterceptor := interceptors.NewArchInterceptor(archUseCase, logger)
	archHandler := handlers.NewArchServiceServer(archInterceptor, logger)
	return &App{
		db:              db,
		grpcServer:      grpcServer,
		logger:          logger,
		archRepository:  archRepository,
		archUseCase:     archUseCase,
		archInterceptor: archInterceptor,
		archHandler:     archHandler,
	}
}
func (a *App) Start(ctx context.Context) error {
	a.grpcServer.AddHandler(&examplepb.ArchService_ServiceDesc, a.archHandler)
	return nil
}
func (a *App) Stop(ctx context.Context) error {
	return nil
}
