package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
)

type ArticleServiceServer struct {
	examplepb.UnimplementedArticleServiceServer
	articleUseCase articleUseCase
	logger         logger
}

func NewArticleServiceServer(articleUseCase articleUseCase, logger logger) *ArticleServiceServer {
	return &ArticleServiceServer{articleUseCase: articleUseCase, logger: logger}
}

func (s *ArticleServiceServer) Create(
	ctx context.Context,
	input *examplepb.ArticleCreate,
) (*examplepb.Article, error) {
	item, err := s.articleUseCase.Create(ctx, encodeArticleCreate(input))
	if err != nil {
		return nil, err
	}
	return decodeArticle(item), nil
}

func (s *ArticleServiceServer) Get(
	ctx context.Context,
	input *examplepb.ArticleGet,
) (*examplepb.Article, error) {
	item, err := s.articleUseCase.Get(ctx, uuid.MustParse(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeArticle(item), nil
}

func (s *ArticleServiceServer) List(
	ctx context.Context,
	filter *examplepb.ArticleFilter,
) (*examplepb.ListArticle, error) {
	items, count, err := s.articleUseCase.List(ctx, encodeArticleFilter(filter))
	if err != nil {
		return nil, err
	}
	return decodeListArticle(items, count), nil
}

func (s *ArticleServiceServer) Update(
	ctx context.Context,
	input *examplepb.ArticleUpdate,
) (*examplepb.Article, error) {
	item, err := s.articleUseCase.Update(ctx, encodeArticleUpdate(input))
	if err != nil {
		return nil, err
	}
	return decodeArticle(item), nil
}

func (s *ArticleServiceServer) Delete(
	ctx context.Context,
	input *examplepb.ArticleDelete,
) (*examplepb.Article, error) {
	article, err := s.articleUseCase.Delete(ctx, encodeArticleDelete(input))
	if err != nil {
		return nil, err
	}
	return decodeArticle(article), nil
}
func (s *ArticleServiceServer) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.ArticleService_ServiceDesc, s)
	return nil
}
