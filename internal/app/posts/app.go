package posts

import (
	"github.com/jmoiron/sqlx"
	likeGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/grpc/like"
	postGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/grpc/post"
	tagGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/grpc/tag"
	likeHttpHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/http/like"
	postHttpHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/http/post"
	tagHttpHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/http/tag"
	likeKafkaHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/kafka/like"
	postKafkaHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/kafka/post"
	tagKafkaHandlers "github.com/mikalai-mitsin/example/internal/app/posts/handlers/kafka/tag"
	likeEvents "github.com/mikalai-mitsin/example/internal/app/posts/repositories/kafka/like"
	postEvents "github.com/mikalai-mitsin/example/internal/app/posts/repositories/kafka/post"
	tagEvents "github.com/mikalai-mitsin/example/internal/app/posts/repositories/kafka/tag"
	likeRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/postgres/like"
	postRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/postgres/post"
	tagRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/postgres/tag"
	likeServices "github.com/mikalai-mitsin/example/internal/app/posts/services/like"
	postServices "github.com/mikalai-mitsin/example/internal/app/posts/services/post"
	tagServices "github.com/mikalai-mitsin/example/internal/app/posts/services/tag"
	likeUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/like"
	postUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/post"
	tagUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	readDB            *sqlx.DB
	writeDB           *sqlx.DB
	logger            *log.Log
	kafkaProducer     *kafka.Producer
	postRepository    *postRepositories.PostRepository
	postService       *postServices.PostService
	postUseCase       *postUseCases.PostUseCase
	httpPostHandler   *postHttpHandlers.PostHandler
	postEventProducer *postEvents.PostEventProducer
	kafkaPostHandler  *postKafkaHandlers.PostHandler
	grpcPostHandler   *postGrpcHandlers.PostServiceServer
	tagRepository     *tagRepositories.TagRepository
	tagService        *tagServices.TagService
	tagUseCase        *tagUseCases.TagUseCase
	httpTagHandler    *tagHttpHandlers.TagHandler
	tagEventProducer  *tagEvents.TagEventProducer
	kafkaTagHandler   *tagKafkaHandlers.TagHandler
	grpcTagHandler    *tagGrpcHandlers.TagServiceServer
	likeRepository    *likeRepositories.LikeRepository
	likeService       *likeServices.LikeService
	likeUseCase       *likeUseCases.LikeUseCase
	httpLikeHandler   *likeHttpHandlers.LikeHandler
	likeEventProducer *likeEvents.LikeEventProducer
	kafkaLikeHandler  *likeKafkaHandlers.LikeHandler
	grpcLikeHandler   *likeGrpcHandlers.LikeServiceServer
}

func NewApp(
	readDB, writeDB *sqlx.DB,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv7Generator,
	kafkaProducer *kafka.Producer,
) *App {
	postRepository := postRepositories.NewPostRepository(readDB, writeDB, logger)
	postService := postServices.NewPostService(postRepository, clock, logger, uuidGenerator)
	postEventProducer := postEvents.NewPostEventProducer(kafkaProducer, logger)
	postUseCase := postUseCases.NewPostUseCase(postService, postEventProducer, logger)
	httpPostHandler := postHttpHandlers.NewPostHandler(postUseCase, logger)
	kafkaPostHandler := postKafkaHandlers.NewPostHandler(postUseCase, logger)
	grpcPostHandler := postGrpcHandlers.NewPostServiceServer(postUseCase, logger)
	tagRepository := tagRepositories.NewTagRepository(readDB, writeDB, logger)
	tagService := tagServices.NewTagService(tagRepository, clock, logger, uuidGenerator)
	tagEventProducer := tagEvents.NewTagEventProducer(kafkaProducer, logger)
	tagUseCase := tagUseCases.NewTagUseCase(tagService, tagEventProducer, logger)
	httpTagHandler := tagHttpHandlers.NewTagHandler(tagUseCase, logger)
	kafkaTagHandler := tagKafkaHandlers.NewTagHandler(tagUseCase, logger)
	grpcTagHandler := tagGrpcHandlers.NewTagServiceServer(tagUseCase, logger)
	likeRepository := likeRepositories.NewLikeRepository(readDB, writeDB, logger)
	likeService := likeServices.NewLikeService(likeRepository, clock, logger, uuidGenerator)
	likeEventProducer := likeEvents.NewLikeEventProducer(kafkaProducer, logger)
	likeUseCase := likeUseCases.NewLikeUseCase(likeService, likeEventProducer, logger)
	httpLikeHandler := likeHttpHandlers.NewLikeHandler(likeUseCase, logger)
	kafkaLikeHandler := likeKafkaHandlers.NewLikeHandler(likeUseCase, logger)
	grpcLikeHandler := likeGrpcHandlers.NewLikeServiceServer(likeUseCase, logger)
	return &App{
		readDB:            readDB,
		writeDB:           writeDB,
		logger:            logger,
		kafkaProducer:     kafkaProducer,
		postRepository:    postRepository,
		postService:       postService,
		postUseCase:       postUseCase,
		httpPostHandler:   httpPostHandler,
		postEventProducer: postEventProducer,
		kafkaPostHandler:  kafkaPostHandler,
		grpcPostHandler:   grpcPostHandler,
		tagRepository:     tagRepository,
		tagService:        tagService,
		tagUseCase:        tagUseCase,
		httpTagHandler:    httpTagHandler,
		tagEventProducer:  tagEventProducer,
		kafkaTagHandler:   kafkaTagHandler,
		grpcTagHandler:    grpcTagHandler,
		likeRepository:    likeRepository,
		likeService:       likeService,
		likeUseCase:       likeUseCase,
		httpLikeHandler:   httpLikeHandler,
		likeEventProducer: likeEventProducer,
		kafkaLikeHandler:  kafkaLikeHandler,
		grpcLikeHandler:   grpcLikeHandler,
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
func (a *App) RegisterKafka(consumer *kafka.Consumer) error {
	consumer.AddHandler(
		kafka.Handler{
			Topic:   "example.posts.post.created",
			GroupID: "example.posts.post",
			Handler: a.kafkaPostHandler,
		},
	)
	consumer.AddHandler(
		kafka.Handler{
			Topic:   "example.posts.tag.created",
			GroupID: "example.posts.tag",
			Handler: a.kafkaTagHandler,
		},
	)
	consumer.AddHandler(
		kafka.Handler{
			Topic:   "example.posts.like.created",
			GroupID: "example.posts.like",
			Handler: a.kafkaLikeHandler,
		},
	)
	return nil
}
