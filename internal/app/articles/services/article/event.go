package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
)

type ArticleEventService struct {
	articleEventProducer articleEventProducer
	logger               logger
}

func NewArticleEventService(
	articleEventProducer articleEventProducer,
	logger logger,
) *ArticleEventService {
	return &ArticleEventService{articleEventProducer: articleEventProducer, logger: logger}
}
func (s *ArticleEventService) Send(ctx context.Context, _ dtx.TX, article entities.Article) error {
	if err := s.articleEventProducer.Send(ctx, article); err != nil {
		return err
	}
	return nil
}
