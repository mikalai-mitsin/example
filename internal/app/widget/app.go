package widget

import (
	"github.com/jmoiron/sqlx"
	handlers "github.com/mikalai-mitsin/example/internal/app/widget/handlers/grpc"
	"github.com/mikalai-mitsin/example/internal/app/widget/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/widget/services"
	"github.com/mikalai-mitsin/example/internal/app/widget/usecases"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db               *sqlx.DB
	logger           *log.Log
	widgetRepository *postgres.WidgetRepository
	widgetService    *services.WidgetService
	widgetUseCase    *usecases.WidgetUseCase
	widgetHandler    *handlers.WidgetServiceServer
}

func NewApp(
	db *sqlx.DB,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv4Generator,
) *App {
	widgetRepository := postgres.NewWidgetRepository(db, logger)
	widgetService := services.NewWidgetService(widgetRepository, clock, logger, uuidGenerator)
	widgetUseCase := usecases.NewWidgetUseCase(widgetService, logger)
	widgetHandler := handlers.NewWidgetServiceServer(widgetUseCase, logger)
	return &App{
		db:               db,
		logger:           logger,
		widgetRepository: widgetRepository,
		widgetService:    widgetService,
		widgetUseCase:    widgetUseCase,
		widgetHandler:    widgetHandler,
	}
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.WidgetService_ServiceDesc, a.widgetRepository)
	return nil
}
