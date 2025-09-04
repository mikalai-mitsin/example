package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
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

func (u *LikeService) Create(
	ctx context.Context,
	tx dtx.TX,
	create entities.LikeCreate,
) (entities.Like, error) {
	if err := create.Validate(); err != nil {
		return entities.Like{}, err
	}
	now := u.clock.Now().UTC()
	like := entities.Like{
		ID:        u.uuid.NewUUID(),
		UpdatedAt: now,
		CreatedAt: now,
		PostId:    create.PostId,
		Value:     create.Value,
		UserId:    create.UserId,
	}
	if err := u.likeRepository.Create(ctx, tx, like); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
func (u *LikeService) Get(ctx context.Context, id uuid.UUID) (entities.Like, error) {
	like, err := u.likeRepository.Get(ctx, id)
	if err != nil {
		return entities.Like{}, err
	}
	return like, nil
}

func (u *LikeService) List(
	ctx context.Context,
	filter entities.LikeFilter,
) ([]entities.Like, uint64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	like, err := u.likeRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.likeRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return like, count, nil
}

func (u *LikeService) Update(
	ctx context.Context,
	tx dtx.TX,
	update entities.LikeUpdate,
) (entities.Like, error) {
	if err := update.Validate(); err != nil {
		return entities.Like{}, err
	}
	like, err := u.likeRepository.Get(ctx, update.ID)
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
	like.UpdatedAt = u.clock.Now().UTC()
	if err := u.likeRepository.Update(ctx, tx, like); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
func (u *LikeService) Delete(ctx context.Context, tx dtx.TX, id uuid.UUID) error {
	if err := u.likeRepository.Delete(ctx, tx, id); err != nil {
		return err
	}
	return nil
}
