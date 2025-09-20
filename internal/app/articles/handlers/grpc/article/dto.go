package handlers

import (
	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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
		IsDeleted:  nil,
		OrderBy:    []entities.ArticleOrdering{},
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = pointer.Of(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = pointer.Of(input.GetPageNumber().GetValue())
	}
	if input.GetIsDeleted() != nil {
		filter.IsDeleted = pointer.Of(input.GetIsDeleted().GetValue())
	}
	if input.GetSearch() != nil {
		filter.Search = pointer.Of(input.GetSearch().GetValue())
	}
	for _, orderBy := range input.GetOrderBy() {
		filter.OrderBy = append(filter.OrderBy, entities.ArticleOrdering(orderBy))
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
		update.IsPublished = pointer.Of(bool(input.GetIsPublished().GetValue()))
	}
	return update
}
func decodeArticle(article entities.Article) *examplepb.Article {
	response := &examplepb.Article{
		Id:          article.ID.String(),
		CreatedAt:   timestamppb.New(article.CreatedAt),
		UpdatedAt:   timestamppb.New(article.UpdatedAt),
		DeletedAt:   nil,
		Title:       article.Title,
		Subtitle:    article.Subtitle,
		Body:        article.Body,
		IsPublished: article.IsPublished,
	}
	if article.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*article.DeletedAt)
	}
	return response
}
func decodeListArticle(items []entities.Article, count uint64) *examplepb.ListArticle {
	response := &examplepb.ListArticle{
		Items: make([]*examplepb.Article, 0, len(items)),
		Count: count,
	}
	for _, article := range items {
		response.Items = append(response.Items, decodeArticle(article))
	}
	return response
}
func decodeArticleUpdate(update entities.ArticleUpdate) *examplepb.ArticleUpdate {
	result := &examplepb.ArticleUpdate{
		Id:          string(update.ID.String()),
		Title:       wrapperspb.String(*update.Title),
		Subtitle:    wrapperspb.String(*update.Subtitle),
		Body:        wrapperspb.String(*update.Body),
		IsPublished: wrapperspb.Bool(*update.IsPublished),
	}
	return result
}
