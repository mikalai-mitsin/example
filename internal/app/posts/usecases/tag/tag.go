package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type TagUseCase struct {
	tagService      tagService
	tagEventService tagEventService
	dtxManager      dtxManager
	logger          logger
}

func NewTagUseCase(
	tagService tagService,
	tagEventService tagEventService,
	dtxManager dtxManager,
	logger logger,
) *TagUseCase {
	return &TagUseCase{
		tagService:      tagService,
		tagEventService: tagEventService,
		dtxManager:      dtxManager,
		logger:          logger,
	}
}
func (u *TagUseCase) Create(ctx context.Context, create entities.TagCreate) (entities.Tag, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	tag, err := u.tagService.Create(ctx, tx, create)
	if err != nil {
		return entities.Tag{}, err
	}
	if err := u.tagEventService.Send(ctx, tx, tag); err != nil {
		return entities.Tag{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}
func (u *TagUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Tag, error) {
	tag, err := u.tagService.Get(ctx, id)
	if err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}

func (u *TagUseCase) List(
	ctx context.Context,
	filter entities.TagFilter,
) ([]entities.Tag, uint64, error) {
	tags, count, err := u.tagService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return tags, count, nil
}
func (u *TagUseCase) Update(ctx context.Context, update entities.TagUpdate) (entities.Tag, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	tag, err := u.tagService.Update(ctx, tx, update)
	if err != nil {
		return entities.Tag{}, err
	}
	if err := u.tagEventService.Send(ctx, tx, tag); err != nil {
		return entities.Tag{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}
func (u *TagUseCase) Delete(ctx context.Context, id uuid.UUID) (entities.Tag, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	tag, err := u.tagService.Delete(ctx, tx, id)
	if err != nil {
		return entities.Tag{}, err
	}
	if err := u.tagEventService.Send(ctx, tx, tag); err != nil {
		return entities.Tag{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}
