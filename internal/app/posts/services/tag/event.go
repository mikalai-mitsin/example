package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type TagEventService struct {
	tagEventProducer tagEventProducer
	logger           logger
}

func NewTagEventService(tagEventProducer tagEventProducer, logger logger) *TagEventService {
	return &TagEventService{tagEventProducer: tagEventProducer, logger: logger}
}
func (s *TagEventService) Created(ctx context.Context, _ dtx.TX, tag entities.Tag) error {
	if err := s.tagEventProducer.Created(ctx, tag); err != nil {
		return err
	}
	return nil
}
func (s *TagEventService) Updated(ctx context.Context, _ dtx.TX, tag entities.Tag) error {
	if err := s.tagEventProducer.Updated(ctx, tag); err != nil {
		return err
	}
	return nil
}
func (s *TagEventService) Deleted(ctx context.Context, _ dtx.TX, id uuid.UUID) error {
	if err := s.tagEventProducer.Deleted(ctx, id); err != nil {
		return err
	}
	return nil
}
