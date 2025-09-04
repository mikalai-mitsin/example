package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type TagService struct {
	tagRepository tagRepository
	clock         clock
	logger        logger
	uuid          uuidGenerator
}

func NewTagService(
	tagRepository tagRepository,
	clock clock,
	logger logger,
	uuid uuidGenerator,
) *TagService {
	return &TagService{tagRepository: tagRepository, clock: clock, logger: logger, uuid: uuid}
}

func (u *TagService) Create(
	ctx context.Context,
	tx dtx.TX,
	create entities.TagCreate,
) (entities.Tag, error) {
	if err := create.Validate(); err != nil {
		return entities.Tag{}, err
	}
	now := u.clock.Now().UTC()
	tag := entities.Tag{
		ID:        u.uuid.NewUUID(),
		UpdatedAt: now,
		CreatedAt: now,
		PostId:    create.PostId,
		Value:     create.Value,
	}
	if err := u.tagRepository.Create(ctx, tx, tag); err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}
func (u *TagService) Get(ctx context.Context, id uuid.UUID) (entities.Tag, error) {
	tag, err := u.tagRepository.Get(ctx, id)
	if err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}

func (u *TagService) List(
	ctx context.Context,
	filter entities.TagFilter,
) ([]entities.Tag, uint64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	tag, err := u.tagRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.tagRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return tag, count, nil
}

func (u *TagService) Update(
	ctx context.Context,
	tx dtx.TX,
	update entities.TagUpdate,
) (entities.Tag, error) {
	if err := update.Validate(); err != nil {
		return entities.Tag{}, err
	}
	tag, err := u.tagRepository.Get(ctx, update.ID)
	if err != nil {
		return entities.Tag{}, err
	}
	{
		if update.PostId != nil {
			tag.PostId = *update.PostId
		}
		if update.Value != nil {
			tag.Value = *update.Value
		}
	}
	tag.UpdatedAt = u.clock.Now().UTC()
	if err := u.tagRepository.Update(ctx, tx, tag); err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}
func (u *TagService) Delete(ctx context.Context, tx dtx.TX, id uuid.UUID) error {
	if err := u.tagRepository.Delete(ctx, tx, id); err != nil {
		return err
	}
	return nil
}
