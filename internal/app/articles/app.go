package articles

import (
	"github.com/jmoiron/sqlx"
	articleGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/articles/handlers/grpc/article"
	articleHttpHandlers "github.com/mikalai-mitsin/example/internal/app/articles/handlers/http/article"
	articleKafkaHandlers "github.com/mikalai-mitsin/example/internal/app/articles/handlers/kafka/article"
	articleKafkaRepositories "github.com/mikalai-mitsin/example/internal/app/articles/repositories/kafka/article"
	articlePostgresRepositories "github.com/mikalai-mitsin/example/internal/app/articles/repositories/postgres/article"
	articleServices "github.com/mikalai-mitsin/example/internal/app/articles/services/article"
	articleUseCases "github.com/mikalai-mitsin/example/internal/app/articles/usecases/article"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type App struct {
	readDB               *sqlx.DB
	writeDB              *sqlx.DB
	dtxManager           *dtx.Manager
	logger               log.Logger
	kafkaProducer        *kafka.Producer
	articleRepository    *articlePostgresRepositories.ArticleRepository
	articleService       *articleServices.ArticleService
	articleUseCase       *articleUseCases.ArticleUseCase
	httpArticleHandler   *articleHttpHandlers.ArticleHandler
	articleEventProducer *articleKafkaRepositories.ArticleEventProducer
	articleEventService  *articleServices.ArticleEventService
	kafkaArticleHandler  *articleKafkaHandlers.ArticleHandler
	grpcArticleHandler   *articleGrpcHandlers.ArticleServiceServer
}

func NewApp(
	readDB, writeDB *sqlx.DB,
	dtxManager *dtx.Manager,
	logger log.Logger,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv7Generator,
	kafkaProducer *kafka.Producer,
) *App {
	articleRepository := articlePostgresRepositories.NewArticleRepository(readDB, writeDB, logger)
	articleService := articleServices.NewArticleService(
		articleRepository,
		clock,
		logger,
		uuidGenerator,
	)
	articleEventProducer := articleKafkaRepositories.NewArticleEventProducer(kafkaProducer, logger)
	articleEventService := articleServices.NewArticleEventService(articleEventProducer, logger)
	articleUseCase := articleUseCases.NewArticleUseCase(
		articleService,
		articleEventService,
		dtxManager,
		logger,
	)
	httpArticleHandler := articleHttpHandlers.NewArticleHandler(articleUseCase, logger)
	kafkaArticleHandler := articleKafkaHandlers.NewArticleHandler(articleUseCase, logger)
	grpcArticleHandler := articleGrpcHandlers.NewArticleServiceServer(articleUseCase, logger)
	return &App{
		readDB:               readDB,
		writeDB:              writeDB,
		dtxManager:           dtxManager,
		logger:               logger,
		kafkaProducer:        kafkaProducer,
		articleRepository:    articleRepository,
		articleService:       articleService,
		articleUseCase:       articleUseCase,
		httpArticleHandler:   httpArticleHandler,
		articleEventProducer: articleEventProducer,
		articleEventService:  articleEventService,
		kafkaArticleHandler:  kafkaArticleHandler,
		grpcArticleHandler:   grpcArticleHandler,
	}
}
func (a *App) RegisterHTTP(httpServer *http.Server) error {
	if err := a.httpArticleHandler.RegisterHTTP(httpServer); err != nil {
		return err
	}
	return nil
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	if err := a.grpcArticleHandler.RegisterGRPC(grpcServer); err != nil {
		return err
	}
	return nil
}
func (a *App) RegisterKafka(consumer *kafka.Consumer) error {
	if err := a.kafkaArticleHandler.RegisterKafka(consumer); err != nil {
		return err
	}
	return nil
}
