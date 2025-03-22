package post

import (
	"github.com/jmoiron/sqlx"
	grpcHandlers "github.com/mikalai-mitsin/example/internal/app/post/handlers/grpc"
	httpHandlers "github.com/mikalai-mitsin/example/internal/app/post/handlers/http"
	"github.com/mikalai-mitsin/example/internal/app/post/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/post/services"
	"github.com/mikalai-mitsin/example/internal/app/post/usecases"
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
	postRepository  *postgres.PostRepository
	postService     *services.PostService
	postUseCase     *usecases.PostUseCase
	httpPostHandler *httpHandlers.PostHandler
	grpcPostHandler *grpcHandlers.PostServiceServer
}

func NewApp(
	db *sqlx.DB,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv4Generator,
) *App {
	postRepository := postgres.NewPostRepository(db, logger)
	postService := services.NewPostService(postRepository, clock, logger, uuidGenerator)
	postUseCase := usecases.NewPostUseCase(postService, logger)
	httpPostHandler := httpHandlers.NewPostHandler(postUseCase, logger)
	grpcPostHandler := grpcHandlers.NewPostServiceServer(postUseCase, logger)
	return &App{
		db:              db,
		logger:          logger,
		postRepository:  postRepository,
		postService:     postService,
		postUseCase:     postUseCase,
		httpPostHandler: httpPostHandler,
		grpcPostHandler: grpcPostHandler,
	}
}
func (a *App) RegisterHTTP(httpServer *http.Server) error {
	httpServer.Mount("/api/v1/posts/", a.httpPostHandler.ChiRouter())
	return nil
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.PostService_ServiceDesc, a.grpcPostHandler)
	return nil
}
