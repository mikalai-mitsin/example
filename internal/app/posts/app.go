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
	likeKafkaRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/kafka/like"
	postKafkaRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/kafka/post"
	tagKafkaRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/kafka/tag"
	likePostgresRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/postgres/like"
	postPostgresRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/postgres/post"
	tagPostgresRepositories "github.com/mikalai-mitsin/example/internal/app/posts/repositories/postgres/tag"
	likeServices "github.com/mikalai-mitsin/example/internal/app/posts/services/like"
	postServices "github.com/mikalai-mitsin/example/internal/app/posts/services/post"
	tagServices "github.com/mikalai-mitsin/example/internal/app/posts/services/tag"
	likeUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/like"
	postUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/post"
	tagUseCases "github.com/mikalai-mitsin/example/internal/app/posts/usecases/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type App struct {
	readDB            *sqlx.DB
	writeDB           *sqlx.DB
	dtxManager        *dtx.Manager
	logger            log.Logger
	kafkaProducer     *kafka.Producer
	postRepository    *postPostgresRepositories.PostRepository
	postService       *postServices.PostService
	postUseCase       *postUseCases.PostUseCase
	httpPostHandler   *postHttpHandlers.PostHandler
	postEventProducer *postKafkaRepositories.PostEventProducer
	postEventService  *postServices.PostEventService
	kafkaPostHandler  *postKafkaHandlers.PostHandler
	grpcPostHandler   *postGrpcHandlers.PostServiceServer
	tagRepository     *tagPostgresRepositories.TagRepository
	tagService        *tagServices.TagService
	tagUseCase        *tagUseCases.TagUseCase
	httpTagHandler    *tagHttpHandlers.TagHandler
	tagEventProducer  *tagKafkaRepositories.TagEventProducer
	tagEventService   *tagServices.TagEventService
	kafkaTagHandler   *tagKafkaHandlers.TagHandler
	grpcTagHandler    *tagGrpcHandlers.TagServiceServer
	likeRepository    *likePostgresRepositories.LikeRepository
	likeService       *likeServices.LikeService
	likeUseCase       *likeUseCases.LikeUseCase
	httpLikeHandler   *likeHttpHandlers.LikeHandler
	likeEventProducer *likeKafkaRepositories.LikeEventProducer
	likeEventService  *likeServices.LikeEventService
	kafkaLikeHandler  *likeKafkaHandlers.LikeHandler
	grpcLikeHandler   *likeGrpcHandlers.LikeServiceServer
}

func NewApp(
	readDB, writeDB *sqlx.DB,
	dtxManager *dtx.Manager,
	logger log.Logger,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv7Generator,
	kafkaProducer *kafka.Producer,
) *App {
	postRepository := postPostgresRepositories.NewPostRepository(readDB, writeDB, logger)
	postService := postServices.NewPostService(postRepository, clock, logger, uuidGenerator)
	postEventProducer := postKafkaRepositories.NewPostEventProducer(kafkaProducer, logger)
	postEventService := postServices.NewPostEventService(postEventProducer, logger)
	postUseCase := postUseCases.NewPostUseCase(postService, postEventService, dtxManager, logger)
	httpPostHandler := postHttpHandlers.NewPostHandler(postUseCase, logger)
	kafkaPostHandler := postKafkaHandlers.NewPostHandler(postUseCase, logger)
	grpcPostHandler := postGrpcHandlers.NewPostServiceServer(postUseCase, logger)
	tagRepository := tagPostgresRepositories.NewTagRepository(readDB, writeDB, logger)
	tagService := tagServices.NewTagService(tagRepository, clock, logger, uuidGenerator)
	tagEventProducer := tagKafkaRepositories.NewTagEventProducer(kafkaProducer, logger)
	tagEventService := tagServices.NewTagEventService(tagEventProducer, logger)
	tagUseCase := tagUseCases.NewTagUseCase(tagService, tagEventService, dtxManager, logger)
	httpTagHandler := tagHttpHandlers.NewTagHandler(tagUseCase, logger)
	kafkaTagHandler := tagKafkaHandlers.NewTagHandler(tagUseCase, logger)
	grpcTagHandler := tagGrpcHandlers.NewTagServiceServer(tagUseCase, logger)
	likeRepository := likePostgresRepositories.NewLikeRepository(readDB, writeDB, logger)
	likeService := likeServices.NewLikeService(likeRepository, clock, logger, uuidGenerator)
	likeEventProducer := likeKafkaRepositories.NewLikeEventProducer(kafkaProducer, logger)
	likeEventService := likeServices.NewLikeEventService(likeEventProducer, logger)
	likeUseCase := likeUseCases.NewLikeUseCase(likeService, likeEventService, dtxManager, logger)
	httpLikeHandler := likeHttpHandlers.NewLikeHandler(likeUseCase, logger)
	kafkaLikeHandler := likeKafkaHandlers.NewLikeHandler(likeUseCase, logger)
	grpcLikeHandler := likeGrpcHandlers.NewLikeServiceServer(likeUseCase, logger)
	return &App{
		readDB:            readDB,
		writeDB:           writeDB,
		dtxManager:        dtxManager,
		logger:            logger,
		kafkaProducer:     kafkaProducer,
		postRepository:    postRepository,
		postService:       postService,
		postUseCase:       postUseCase,
		httpPostHandler:   httpPostHandler,
		postEventProducer: postEventProducer,
		postEventService:  postEventService,
		kafkaPostHandler:  kafkaPostHandler,
		grpcPostHandler:   grpcPostHandler,
		tagRepository:     tagRepository,
		tagService:        tagService,
		tagUseCase:        tagUseCase,
		httpTagHandler:    httpTagHandler,
		tagEventProducer:  tagEventProducer,
		tagEventService:   tagEventService,
		kafkaTagHandler:   kafkaTagHandler,
		grpcTagHandler:    grpcTagHandler,
		likeRepository:    likeRepository,
		likeService:       likeService,
		likeUseCase:       likeUseCase,
		httpLikeHandler:   httpLikeHandler,
		likeEventProducer: likeEventProducer,
		likeEventService:  likeEventService,
		kafkaLikeHandler:  kafkaLikeHandler,
		grpcLikeHandler:   grpcLikeHandler,
	}
}
func (a *App) RegisterHTTP(httpServer *http.Server) error {
	if err := a.httpPostHandler.RegisterHTTP(httpServer); err != nil {
		return err
	}
	if err := a.httpTagHandler.RegisterHTTP(httpServer); err != nil {
		return err
	}
	if err := a.httpLikeHandler.RegisterHTTP(httpServer); err != nil {
		return err
	}
	return nil
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	if err := a.grpcPostHandler.RegisterGRPC(grpcServer); err != nil {
		return err
	}
	if err := a.grpcTagHandler.RegisterGRPC(grpcServer); err != nil {
		return err
	}
	if err := a.grpcLikeHandler.RegisterGRPC(grpcServer); err != nil {
		return err
	}
	return nil
}
func (a *App) RegisterKafka(consumer *kafka.Consumer) error {
	if err := a.kafkaPostHandler.RegisterKafka(consumer); err != nil {
		return err
	}
	if err := a.kafkaTagHandler.RegisterKafka(consumer); err != nil {
		return err
	}
	if err := a.kafkaLikeHandler.RegisterKafka(consumer); err != nil {
		return err
	}
	return nil
}
