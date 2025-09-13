package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
)

type TagEventService struct {
	tagEventProducer tagEventProducer
	logger           logger
}

func NewTagEventService(tagEventProducer tagEventProducer, logger logger) *TagEventService {
	return &TagEventService{tagEventProducer: tagEventProducer, logger: logger}
}
func (s *TagEventService) Send(ctx context.Context, _ dtx.TX, tag entities.Tag) error {
	if err := s.tagEventProducer.Send(ctx, tag); err != nil {
		return err
	}
	return nil
}
