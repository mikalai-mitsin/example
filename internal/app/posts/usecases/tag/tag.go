package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type TagUseCase struct {
	tagService tagService
	logger     logger
}

func NewTagUseCase(tagService tagService, logger logger) *TagUseCase {
	return &TagUseCase{tagService: tagService, logger: logger}
}
func (i *TagUseCase) Create(ctx context.Context, create entities.TagCreate) (entities.Tag, error) {
	tag, err := i.tagService.Create(ctx, create)
	if err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}
func (i *TagUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Tag, error) {
	tag, err := i.tagService.Get(ctx, id)
	if err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}

func (i *TagUseCase) List(
	ctx context.Context,
	filter entities.TagFilter,
) ([]entities.Tag, uint64, error) {
	items, count, err := i.tagService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *TagUseCase) Update(ctx context.Context, update entities.TagUpdate) (entities.Tag, error) {
	updated, err := i.tagService.Update(ctx, update)
	if err != nil {
		return entities.Tag{}, err
	}
	return updated, nil
}
func (i *TagUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.tagService.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
