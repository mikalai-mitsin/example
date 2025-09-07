package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type LikeUseCase struct {
	likeService      likeService
	likeEventService likeEventService
	dtxManager       dtxManager
	logger           logger
}

func NewLikeUseCase(
	likeService likeService,
	likeEventService likeEventService,
	dtxManager dtxManager,
	logger logger,
) *LikeUseCase {
	return &LikeUseCase{
		likeService:      likeService,
		likeEventService: likeEventService,
		dtxManager:       dtxManager,
		logger:           logger,
	}
}

func (u *LikeUseCase) Create(
	ctx context.Context,
	create entities.LikeCreate,
) (entities.Like, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	like, err := u.likeService.Create(ctx, tx, create)
	if err != nil {
		return entities.Like{}, err
	}
	if err := u.likeEventService.Created(ctx, tx, like); err != nil {
		return entities.Like{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
func (u *LikeUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Like, error) {
	like, err := u.likeService.Get(ctx, id)
	if err != nil {
		return entities.Like{}, err
	}
	return like, nil
}

func (u *LikeUseCase) List(
	ctx context.Context,
	filter entities.LikeFilter,
) ([]entities.Like, uint64, error) {
	likes, count, err := u.likeService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return likes, count, nil
}

func (u *LikeUseCase) Update(
	ctx context.Context,
	update entities.LikeUpdate,
) (entities.Like, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	like, err := u.likeService.Update(ctx, tx, update)
	if err != nil {
		return entities.Like{}, err
	}
	if err := u.likeEventService.Updated(ctx, tx, like); err != nil {
		return entities.Like{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Like{}, err
	}
	return like, nil
}
func (u *LikeUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	if err := u.likeService.Delete(ctx, tx, id); err != nil {
		return err
	}
	if err := u.likeEventService.Deleted(ctx, tx, id); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
