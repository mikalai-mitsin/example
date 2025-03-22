package comment

import (
	"github.com/jmoiron/sqlx"
	grpcHandlers "github.com/mikalai-mitsin/example/internal/app/comment/handlers/grpc"
	httpHandlers "github.com/mikalai-mitsin/example/internal/app/comment/handlers/http"
	"github.com/mikalai-mitsin/example/internal/app/comment/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/comment/services"
	"github.com/mikalai-mitsin/example/internal/app/comment/usecases"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db                 *sqlx.DB
	logger             *log.Log
	commentRepository  *postgres.CommentRepository
	commentService     *services.CommentService
	commentUseCase     *usecases.CommentUseCase
	httpCommentHandler *httpHandlers.CommentHandler
	grpcCommentHandler *grpcHandlers.CommentServiceServer
}

func NewApp(
	db *sqlx.DB,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv4Generator,
) *App {
	commentRepository := postgres.NewCommentRepository(db, logger)
	commentService := services.NewCommentService(commentRepository, clock, logger, uuidGenerator)
	commentUseCase := usecases.NewCommentUseCase(commentService, logger)
	httpCommentHandler := httpHandlers.NewCommentHandler(commentUseCase, logger)
	grpcCommentHandler := grpcHandlers.NewCommentServiceServer(commentUseCase, logger)
	return &App{
		db:                 db,
		logger:             logger,
		commentRepository:  commentRepository,
		commentService:     commentService,
		commentUseCase:     commentUseCase,
		httpCommentHandler: httpCommentHandler,
		grpcCommentHandler: grpcCommentHandler,
	}
}
func (a *App) RegisterHTTP(httpServer *http.Server) error {
	httpServer.Mount("/api/v1/comments/", a.httpCommentHandler.ChiRouter())
	return nil
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.CommentService_ServiceDesc, a.grpcCommentHandler)
	return nil
}
