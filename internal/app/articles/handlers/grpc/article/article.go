package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/i18n"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type ArticleServiceServer struct {
	examplepb.UnimplementedArticleServiceServer
	articleUseCase articleUseCase
	logger         logger
	tr             *i18n.Translator
}

func NewArticleServiceServer(
	articleUseCase articleUseCase,
	logger logger,
	tr *i18n.Translator,
) *ArticleServiceServer {
	return &ArticleServiceServer{articleUseCase: articleUseCase, logger: logger, tr: tr}
}

func (s *ArticleServiceServer) Create(
	ctx context.Context,
	input *examplepb.ArticleCreate,
) (*examplepb.Article, error) {
	item, err := s.articleUseCase.Create(ctx, encodeArticleCreate(input))
	if err != nil {
		return nil, err
	}
	return decodeArticle(item, s.tr), nil
}

func (s *ArticleServiceServer) Get(
	ctx context.Context,
	input *examplepb.ArticleGet,
) (*examplepb.Article, error) {
	item, err := s.articleUseCase.Get(ctx, uuid.MustParse(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeArticle(item, s.tr), nil
}

func (s *ArticleServiceServer) List(
	ctx context.Context,
	filter *examplepb.ArticleFilter,
) (*examplepb.ListArticle, error) {
	items, count, err := s.articleUseCase.List(ctx, encodeArticleFilter(filter))
	if err != nil {
		return nil, err
	}
	return decodeListArticle(items, count, s.tr), nil
}

func (s *ArticleServiceServer) Update(
	ctx context.Context,
	input *examplepb.ArticleUpdate,
) (*examplepb.Article, error) {
	item, err := s.articleUseCase.Update(ctx, encodeArticleUpdate(input))
	if err != nil {
		return nil, err
	}
	return decodeArticle(item, s.tr), nil
}

func (s *ArticleServiceServer) Delete(
	ctx context.Context,
	input *examplepb.ArticleDelete,
) (*examplepb.Article, error) {
	article, err := s.articleUseCase.Delete(ctx, uuid.MustParse(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeArticle(article, s.tr), nil
}
func (s *ArticleServiceServer) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.ArticleService_ServiceDesc, s)
	return nil
}
