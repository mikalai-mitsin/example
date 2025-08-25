package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type LikeUseCase struct {
	likeService       likeService
	likeEventProducer likeEventProducer
	logger            logger
}

func NewLikeUseCase(
	likeService likeService,
	likeEventProducer likeEventProducer,
	logger logger,
) *LikeUseCase {
	return &LikeUseCase{
		likeService:       likeService,
		likeEventProducer: likeEventProducer,
		logger:            logger,
	}
}

func (i *LikeUseCase) Create(
	ctx context.Context,
	create entities.LikeCreate,
) (entities.Like, error) {
	like, err := i.likeService.Create(ctx, create)
	if err != nil {
		return entities.Like{}, err
	}
	if err := i.likeEventProducer.Created(ctx, like); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
func (i *LikeUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Like, error) {
	like, err := i.likeService.Get(ctx, id)
	if err != nil {
		return entities.Like{}, err
	}
	return like, nil
}

func (i *LikeUseCase) List(
	ctx context.Context,
	filter entities.LikeFilter,
) ([]entities.Like, uint64, error) {
	likes, count, err := i.likeService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return likes, count, nil
}

func (i *LikeUseCase) Update(
	ctx context.Context,
	update entities.LikeUpdate,
) (entities.Like, error) {
	like, err := i.likeService.Update(ctx, update)
	if err != nil {
		return entities.Like{}, err
	}
	if err := i.likeEventProducer.Updated(ctx, like); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
func (i *LikeUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.likeService.Delete(ctx, id); err != nil {
		return err
	}
	if err := i.likeEventProducer.Deleted(ctx, id); err != nil {
		return err
	}
	return nil
}
