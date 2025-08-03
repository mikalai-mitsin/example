package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type LikeUseCase struct {
	likeService likeService
	logger      logger
}

func NewLikeUseCase(likeService likeService, logger logger) *LikeUseCase {
	return &LikeUseCase{likeService: likeService, logger: logger}
}

func (i *LikeUseCase) Create(
	ctx context.Context,
	create entities.LikeCreate,
) (entities.Like, error) {
	like, err := i.likeService.Create(ctx, create)
	if err != nil {
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
	items, count, err := i.likeService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (i *LikeUseCase) Update(
	ctx context.Context,
	update entities.LikeUpdate,
) (entities.Like, error) {
	updated, err := i.likeService.Update(ctx, update)
	if err != nil {
		return entities.Like{}, err
	}
	return updated, nil
}
func (i *LikeUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.likeService.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
