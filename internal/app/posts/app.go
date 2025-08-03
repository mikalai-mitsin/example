package posts

import (
	"github.com/jmoiron/sqlx"
	likeGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/grpc/like"
	postGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/grpc/post"
	tagGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/grpc/tag"
	likeHttpHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/http/like"
	postHttpHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/http/post"
	tagHttpHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/http/tag"
	likeRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/like"
	postRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/post"
	tagRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/tag"
	likeServices "github.com/mikalai-mitsin/example/internal/app/posts/services/like"
	postServices "github.com/mikalai-mitsin/example/internal/app/posts/services/post"
	tagServices "github.com/mikalai-mitsin/example/internal/app/posts/services/tag"
	likeUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/like"
	postUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/post"
	tagUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/tag"
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
	postRepository  *postRepositories.PostRepository
	postService     *postServices.PostService
	postUseCase     *postUseCases.PostUseCase
	httpPostHandler *postHttpHandlers.PostHandler
	grpcPostHandler *postGrpcHandlers.PostServiceServer
	tagRepository   *tagRepositories.TagRepository
	tagService      *tagServices.TagService
	tagUseCase      *tagUseCases.TagUseCase
	httpTagHandler  *tagHttpHandlers.TagHandler
	grpcTagHandler  *tagGrpcHandlers.TagServiceServer
	likeRepository  *likeRepositories.LikeRepository
	likeService     *likeServices.LikeService
	likeUseCase     *likeUseCases.LikeUseCase
	httpLikeHandler *likeHttpHandlers.LikeHandler
	grpcLikeHandler *likeGrpcHandlers.LikeServiceServer
}

func NewApp(
	db *sqlx.DB,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv7Generator,
) *App {
	postRepository := postRepositories.NewPostRepository(db, logger)
	postService := postServices.NewPostService(postRepository, clock, logger, uuidGenerator)
	postUseCase := postUseCases.NewPostUseCase(postService, logger)
	httpPostHandler := postHttpHandlers.NewPostHandler(postUseCase, logger)
	grpcPostHandler := postGrpcHandlers.NewPostServiceServer(postUseCase, logger)
	tagRepository := tagRepositories.NewTagRepository(db, logger)
	tagService := tagServices.NewTagService(tagRepository, clock, logger, uuidGenerator)
	tagUseCase := tagUseCases.NewTagUseCase(tagService, logger)
	httpTagHandler := tagHttpHandlers.NewTagHandler(tagUseCase, logger)
	grpcTagHandler := tagGrpcHandlers.NewTagServiceServer(tagUseCase, logger)
	likeRepository := likeRepositories.NewLikeRepository(db, logger)
	likeService := likeServices.NewLikeService(likeRepository, clock, logger, uuidGenerator)
	likeUseCase := likeUseCases.NewLikeUseCase(likeService, logger)
	httpLikeHandler := likeHttpHandlers.NewLikeHandler(likeUseCase, logger)
	grpcLikeHandler := likeGrpcHandlers.NewLikeServiceServer(likeUseCase, logger)
	return &App{
		db:              db,
		logger:          logger,
		postRepository:  postRepository,
		postService:     postService,
		postUseCase:     postUseCase,
		httpPostHandler: httpPostHandler,
		grpcPostHandler: grpcPostHandler,
		tagRepository:   tagRepository,
		tagService:      tagService,
		tagUseCase:      tagUseCase,
		httpTagHandler:  httpTagHandler,
		grpcTagHandler:  grpcTagHandler,
		likeRepository:  likeRepository,
		likeService:     likeService,
		likeUseCase:     likeUseCase,
		httpLikeHandler: httpLikeHandler,
		grpcLikeHandler: grpcLikeHandler,
	}
}
func (a *App) RegisterHTTP(httpServer *http.Server) error {
	httpServer.Mount("/api/v1/posts/posts", a.httpPostHandler.ChiRouter())
	httpServer.Mount("/api/v1/posts/tags", a.httpTagHandler.ChiRouter())
	httpServer.Mount("/api/v1/posts/likes", a.httpLikeHandler.ChiRouter())
	return nil
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.PostService_ServiceDesc, a.grpcPostHandler)
	grpcServer.AddHandler(&examplepb.TagService_ServiceDesc, a.grpcTagHandler)
	grpcServer.AddHandler(&examplepb.LikeService_ServiceDesc, a.grpcLikeHandler)
	return nil
}
