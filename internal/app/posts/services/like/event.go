package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type LikeEventService struct {
	likeEventProducer likeEventProducer
	logger            logger
}

func NewLikeEventService(likeEventProducer likeEventProducer, logger logger) *LikeEventService {
	return &LikeEventService{likeEventProducer: likeEventProducer, logger: logger}
}
func (s *LikeEventService) Created(ctx context.Context, _ dtx.TX, like entities.Like) error {
	if err := s.likeEventProducer.Created(ctx, like); err != nil {
		return err
	}
	return nil
}
func (s *LikeEventService) Updated(ctx context.Context, _ dtx.TX, like entities.Like) error {
	if err := s.likeEventProducer.Updated(ctx, like); err != nil {
		return err
	}
	return nil
}
func (s *LikeEventService) Deleted(ctx context.Context, _ dtx.TX, id uuid.UUID) error {
	if err := s.likeEventProducer.Deleted(ctx, id); err != nil {
		return err
	}
	return nil
}
