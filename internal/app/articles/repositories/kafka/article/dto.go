package repositories

import (
	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	if article.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*article.DeletedAt)
	}
	return response
}
