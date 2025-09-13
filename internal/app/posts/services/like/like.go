package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type LikeService struct {
	likeRepository likeRepository
	clock          clock
	logger         logger
	uuid           uuidGenerator
}

func NewLikeService(
	likeRepository likeRepository,
	clock clock,
	logger logger,
	uuid uuidGenerator,
) *LikeService {
	return &LikeService{likeRepository: likeRepository, clock: clock, logger: logger, uuid: uuid}
}

func (s *LikeService) Create(
	ctx context.Context,
	tx dtx.TX,
	create entities.LikeCreate,
) (entities.Like, error) {
	if err := create.Validate(); err != nil {
		return entities.Like{}, err
	}
	now := s.clock.Now().UTC()
	like := entities.Like{
		ID:        s.uuid.NewUUID(),
		UpdatedAt: now,
		CreatedAt: now,
		PostId:    create.PostId,
		Value:     create.Value,
		UserId:    create.UserId,
	}
	if err := s.likeRepository.Create(ctx, tx, like); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
func (s *LikeService) Get(ctx context.Context, id uuid.UUID) (entities.Like, error) {
	like, err := s.likeRepository.Get(ctx, id)
	if err != nil {
		return entities.Like{}, err
	}
	return like, nil
}

func (s *LikeService) List(
	ctx context.Context,
	filter entities.LikeFilter,
) ([]entities.Like, uint64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	like, err := s.likeRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.likeRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return like, count, nil
}

func (s *LikeService) Update(
	ctx context.Context,
	tx dtx.TX,
	update entities.LikeUpdate,
) (entities.Like, error) {
	if err := update.Validate(); err != nil {
		return entities.Like{}, err
	}
	like, err := s.likeRepository.Get(ctx, update.ID)
	if err != nil {
		return entities.Like{}, err
	}
	{
		if update.PostId != nil {
			like.PostId = *update.PostId
		}
		if update.Value != nil {
			like.Value = *update.Value
		}
		if update.UserId != nil {
			like.UserId = *update.UserId
		}
	}
	like.UpdatedAt = s.clock.Now().UTC()
	if err := s.likeRepository.Update(ctx, tx, like); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
func (s *LikeService) Delete(ctx context.Context, tx dtx.TX, id uuid.UUID) (entities.Like, error) {
	like, err := s.likeRepository.Get(ctx, id)
	if err != nil {
		return entities.Like{}, err
	}
	like.DeletedAt = pointer.Of(s.clock.Now().UTC())
	if err := s.likeRepository.Update(ctx, tx, like); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
