package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
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

func (s *TagService) Create(
	ctx context.Context,
	tx dtx.TX,
	create entities.TagCreate,
) (entities.Tag, error) {
	if err := create.Validate(); err != nil {
		return entities.Tag{}, err
	}
	now := s.clock.Now().UTC()
	tag := entities.Tag{
		ID:        s.uuid.NewUUID(),
		UpdatedAt: now,
		CreatedAt: now,
		PostId:    create.PostId,
		Value:     create.Value,
	}
	if err := s.tagRepository.Create(ctx, tx, tag); err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}
func (s *TagService) Get(ctx context.Context, id uuid.UUID) (entities.Tag, error) {
	tag, err := s.tagRepository.Get(ctx, id)
	if err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}

func (s *TagService) List(
	ctx context.Context,
	filter entities.TagFilter,
) ([]entities.Tag, uint64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	tag, err := s.tagRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.tagRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return tag, count, nil
}

func (s *TagService) Update(
	ctx context.Context,
	tx dtx.TX,
	update entities.TagUpdate,
) (entities.Tag, error) {
	if err := update.Validate(); err != nil {
		return entities.Tag{}, err
	}
	tag, err := s.tagRepository.Get(ctx, update.ID)
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
	tag.UpdatedAt = s.clock.Now().UTC()
	if err := s.tagRepository.Update(ctx, tx, tag); err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}

func (s *TagService) Delete(
	ctx context.Context,
	tx dtx.TX,
	del entities.TagDelete,
) (entities.Tag, error) {
	tag, err := s.tagRepository.Get(ctx, del.ID)
	if err != nil {
		return entities.Tag{}, err
	}
	tag.DeletedAt = pointer.Of(s.clock.Now().UTC())
	if err := s.tagRepository.Update(ctx, tx, tag); err != nil {
		return entities.Tag{}, err
	}
	return tag, nil
}
