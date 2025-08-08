package articles

import (
	"github.com/jmoiron/sqlx"
	articleGrpcHandlers "github.com/mikalai-mitsin/example/internal/app/articles/handlers/grpc/article"
	articleHttpHandlers "github.com/mikalai-mitsin/example/internal/app/articles/handlers/http/article"
	articleRepositories "github.com/mikalai-mitsin/example/internal/app/articles/repositories/article"
	articleServices "github.com/mikalai-mitsin/example/internal/app/articles/services/article"
	articleUseCases "github.com/mikalai-mitsin/example/internal/app/articles/usecases/article"
	"github.com/mikalai-mitsin/example/internal/pkg/clock"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type App struct {
	readDB             *sqlx.DB
	writeDB            *sqlx.DB
	logger             *log.Log
	articleRepository  *articleRepositories.ArticleRepository
	articleService     *articleServices.ArticleService
	articleUseCase     *articleUseCases.ArticleUseCase
	httpArticleHandler *articleHttpHandlers.ArticleHandler
	grpcArticleHandler *articleGrpcHandlers.ArticleServiceServer
}

func NewApp(
	readDB, writeDB *sqlx.DB,
	logger *log.Log,
	clock *clock.Clock,
	uuidGenerator *uuid.UUIDv7Generator,
) *App {
	articleRepository := articleRepositories.NewArticleRepository(readDB, writeDB, logger)
	articleService := articleServices.NewArticleService(
		articleRepository,
		clock,
		logger,
		uuidGenerator,
	)
	articleUseCase := articleUseCases.NewArticleUseCase(articleService, logger)
	httpArticleHandler := articleHttpHandlers.NewArticleHandler(articleUseCase, logger)
	grpcArticleHandler := articleGrpcHandlers.NewArticleServiceServer(articleUseCase, logger)
	return &App{
		readDB:             readDB,
		writeDB:            writeDB,
		logger:             logger,
		articleRepository:  articleRepository,
		articleService:     articleService,
		articleUseCase:     articleUseCase,
		httpArticleHandler: httpArticleHandler,
		grpcArticleHandler: grpcArticleHandler,
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
