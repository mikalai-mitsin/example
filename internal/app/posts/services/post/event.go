package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
)

type PostEventService struct {
	postEventProducer postEventProducer
	logger            logger
}

func NewPostEventService(postEventProducer postEventProducer, logger logger) *PostEventService {
	return &PostEventService{postEventProducer: postEventProducer, logger: logger}
}
func (s *PostEventService) Send(ctx context.Context, _ dtx.TX, post entities.Post) error {
	if err := s.postEventProducer.Send(ctx, post); err != nil {
		return err
	}
	return nil
}
