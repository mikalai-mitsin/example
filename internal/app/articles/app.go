package articles

import (
	"github.com/jmoiron/sqlx"
	articleGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/articles/handlers/grpc/article"
	articleHttpHandlers "github.com/mikalai-mitsin/example/internal/app/articles/handlers/http/article"
	articleKafkaHandlers "github.com/mikalai-mitsin/example/internal/app/articles/handlers/kafka/article"
	articleEvents "github.com/mikalai-mitsin/example/internal/app/articles/repositories/kafka/article"
	articleRepositories "github.com/mikalai-mitsin/example/internal/app/articles/repositories/postgres/article"
	articleServices "github.com/mikalai-mitsin/example/internal/app/articles/services/article"
	articleUseCases "github.com/mikalai-mitsin/example/internal/app/articles/usecases/article"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	readDB               *sqlx.DB
	writeDB              *sqlx.DB
	logger               log.Logger
	kafkaProducer        *kafka.Producer
	articleRepository    *articleRepositories.ArticleRepository
	articleService       *articleServices.ArticleService
	articleUseCase       *articleUseCases.ArticleUseCase
	httpArticleHandler   *articleHttpHandlers.ArticleHandler
	articleEventProducer *articleEvents.ArticleEventProducer
	kafkaArticleHandler  *articleKafkaHandlers.ArticleHandler
	grpcArticleHandler   *articleGrpcHandlers.ArticleServiceServer
}

func NewApp(
	readDB, writeDB *sqlx.DB,
	logger log.Logger,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv7Generator,
	kafkaProducer *kafka.Producer,
) *App {
	articleRepository := articleRepositories.NewArticleRepository(readDB, writeDB, logger)
	articleService := articleServices.NewArticleService(
		articleRepository,
		clock,
		logger,
		uuidGenerator,
	)
	articleEventProducer := articleEvents.NewArticleEventProducer(kafkaProducer, logger)
	articleUseCase := articleUseCases.NewArticleUseCase(
		articleService,
		articleEventProducer,
		logger,
	)
	httpArticleHandler := articleHttpHandlers.NewArticleHandler(articleUseCase, logger)
	kafkaArticleHandler := articleKafkaHandlers.NewArticleHandler(articleUseCase, logger)
	grpcArticleHandler := articleGrpcHandlers.NewArticleServiceServer(articleUseCase, logger)
	return &App{
		readDB:               readDB,
		writeDB:              writeDB,
		logger:               logger,
		kafkaProducer:        kafkaProducer,
		articleRepository:    articleRepository,
		articleService:       articleService,
		articleUseCase:       articleUseCase,
		httpArticleHandler:   httpArticleHandler,
		articleEventProducer: articleEventProducer,
		kafkaArticleHandler:  kafkaArticleHandler,
		grpcArticleHandler:   grpcArticleHandler,
	}
}
func (a *App) RegisterHTTP(httpServer *http.Server) error {
	httpServer.Mount("/api/v1/articles/articles", a.httpArticleHandler.ChiRouter())
	return nil
}
func (a *App) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.ArticleService_ServiceDesc, a.grpcArticleHandler)
	return nil
}
func (a *App) RegisterKafka(consumer *kafka.Consumer) error {
	consumer.AddHandler(
		kafka.NewHandler(
			"example.articles.article.created",
			"example.articles.article.created",
			a.kafkaArticleHandler.Created,
		),
	)
	consumer.AddHandler(
		kafka.NewHandler(
			"example.articles.article.updated",
			"example.articles.article.updated",
			a.kafkaArticleHandler.Updated,
		),
	)
	consumer.AddHandler(
		kafka.NewHandler(
			"example.articles.article.deleted",
			"example.articles.article.deleted",
			a.kafkaArticleHandler.Deleted,
		),
	)
	return nil
}
