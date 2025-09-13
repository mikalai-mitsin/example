package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
)

type LikeEventService struct {
	likeEventProducer likeEventProducer
	logger            logger
}

func NewLikeEventService(likeEventProducer likeEventProducer, logger logger) *LikeEventService {
	return &LikeEventService{likeEventProducer: likeEventProducer, logger: logger}
}
func (s *LikeEventService) Send(ctx context.Context, _ dtx.TX, like entities.Like) error {
	if err := s.likeEventProducer.Send(ctx, like); err != nil {
		return err
	}
	return nil
}
