package post

import (
	"github.com/jmoiron/sqlx"
	handlers "github.com/mikalai-mitsin/example/internal/app/post/handlers/grpc"
	"github.com/mikalai-mitsin/example/internal/app/post/repositories/postgres"
	"github.com/mikalai-mitsin/example/internal/app/post/services"
	"github.com/mikalai-mitsin/example/internal/app/post/usecases"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	db             *sqlx.DB
	logger         *log.Log
	postRepository *postgres.PostRepository
	postService    *services.PostService
	postUseCase    *usecases.PostUseCase
	postHandler    *handlers.PostServiceServer
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
	postHandler := handlers.NewPostServiceServer(postUseCase, logger)
	return &App{
		db:             db,
		logger:         logger,
		postRepository: postRepository,
		postService:    postService,
		postUseCase:    postUseCase,
		postHandler:    postHandler,
	}
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.PostService_ServiceDesc, a.postRepository)
	return nil
}
