package handlers

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
) (*emptypb.Empty, error) {
	if err := s.articleUseCase.Delete(ctx, uuid.MustParse(input.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func encodeArticleCreate(input *examplepb.ArticleCreate) entities.ArticleCreate {
	create := entities.ArticleCreate{
		Title:       input.GetTitle(),
		Subtitle:    input.GetSubtitle(),
		Body:        input.GetBody(),
		IsPublished: input.GetIsPublished(),
	}
	return create
}
func encodeArticleFilter(input *examplepb.ArticleFilter) entities.ArticleFilter {
	filter := entities.ArticleFilter{
		PageSize:   nil,
		PageNumber: nil,
		OrderBy:    input.GetOrderBy(),
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = pointer.Of(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = pointer.Of(input.GetPageNumber().GetValue())
	}
	return filter
}
func encodeArticleUpdate(input *examplepb.ArticleUpdate) entities.ArticleUpdate {
	update := entities.ArticleUpdate{ID: uuid.MustParse(input.GetId())}
	if input.GetTitle() != nil {
		update.Title = pointer.Of(string(input.GetTitle().GetValue()))
	}
	if input.GetSubtitle() != nil {
		update.Subtitle = pointer.Of(string(input.GetSubtitle().GetValue()))
	}
	if input.GetBody() != nil {
		update.Body = pointer.Of(string(input.GetBody().GetValue()))
	}
	if input.GetIsPublished() != nil {
		update.IsPublished = pointer.Of(string(input.GetIsPublished().GetValue()))
	}
	return update
}
func decodeArticle(item entities.Article) *examplepb.Article {
	response := &examplepb.Article{
		Id:          item.ID.String(),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
		Title:       item.Title,
		Subtitle:    item.Subtitle,
		Body:        item.Body,
		IsPublished: item.IsPublished,
	}
	return response
}
func decodeListArticle(items []entities.Article, count uint64) *examplepb.ListArticle {
	response := &examplepb.ListArticle{
		Items: make([]*examplepb.Article, 0, len(items)),
		Count: count,
	}
	for _, item := range items {
		response.Items = append(response.Items, decodeArticle(item))
	}
	return response
}
func decodeArticleUpdate(update entities.ArticleUpdate) *examplepb.ArticleUpdate {
	result := &examplepb.ArticleUpdate{
		Id:          string(update.ID.String()),
		Title:       wrapperspb.String(*update.Title),
		Subtitle:    wrapperspb.String(*update.Subtitle),
		Body:        wrapperspb.String(*update.Body),
		IsPublished: wrapperspb.String(*update.IsPublished),
	}
	return result
}
