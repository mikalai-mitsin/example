package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
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

func (s *ArticleEventService) Created(
	ctx context.Context,
	_ dtx.TX,
	article entities.Article,
) error {
	if err := s.articleEventProducer.Created(ctx, article); err != nil {
		return err
	}
	return nil
}

func (s *ArticleEventService) Updated(
	ctx context.Context,
	_ dtx.TX,
	article entities.Article,
) error {
	if err := s.articleEventProducer.Updated(ctx, article); err != nil {
		return err
	}
	return nil
}
func (s *ArticleEventService) Deleted(ctx context.Context, _ dtx.TX, id uuid.UUID) error {
	if err := s.articleEventProducer.Deleted(ctx, id); err != nil {
		return err
	}
	return nil
}
