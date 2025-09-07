package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostEventService struct {
	postEventProducer postEventProducer
	logger            logger
}

func NewPostEventService(postEventProducer postEventProducer, logger logger) *PostEventService {
	return &PostEventService{postEventProducer: postEventProducer, logger: logger}
}
func (s *PostEventService) Created(ctx context.Context, _ dtx.TX, post entities.Post) error {
	if err := s.postEventProducer.Created(ctx, post); err != nil {
		return err
	}
	return nil
}
func (s *PostEventService) Updated(ctx context.Context, _ dtx.TX, post entities.Post) error {
	if err := s.postEventProducer.Updated(ctx, post); err != nil {
		return err
	}
	return nil
}
func (s *PostEventService) Deleted(ctx context.Context, _ dtx.TX, id uuid.UUID) error {
	if err := s.postEventProducer.Deleted(ctx, id); err != nil {
		return err
	}
	return nil
}
