package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostService struct {
	postRepository postRepository
	clock          clock
	logger         logger
	uuid           uuidGenerator
}

func NewPostService(
	postRepository postRepository,
	clock clock,
	logger logger,
	uuid uuidGenerator,
) *PostService {
	return &PostService{postRepository: postRepository, clock: clock, logger: logger, uuid: uuid}
}

func (s *PostService) Create(
	ctx context.Context,
	tx dtx.TX,
	create entities.PostCreate,
) (entities.Post, error) {
	if err := create.Validate(); err != nil {
		return entities.Post{}, err
	}
	now := s.clock.Now().UTC()
	post := entities.Post{ID: s.uuid.NewUUID(), UpdatedAt: now, CreatedAt: now, Body: create.Body}
	if err := s.postRepository.Create(ctx, tx, post); err != nil {
		return entities.Post{}, err
	}
	return post, nil
}
func (s *PostService) Get(ctx context.Context, id uuid.UUID) (entities.Post, error) {
	post, err := s.postRepository.Get(ctx, id)
	if err != nil {
		return entities.Post{}, err
	}
	return post, nil
}

func (s *PostService) List(
	ctx context.Context,
	filter entities.PostFilter,
) ([]entities.Post, uint64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	post, err := s.postRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.postRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return post, count, nil
}

func (s *PostService) Update(
	ctx context.Context,
	tx dtx.TX,
	update entities.PostUpdate,
) (entities.Post, error) {
	if err := update.Validate(); err != nil {
		return entities.Post{}, err
	}
	post, err := s.postRepository.Get(ctx, update.ID)
	if err != nil {
		return entities.Post{}, err
	}
	{
		if update.Body != nil {
			post.Body = *update.Body
		}
	}
	post.UpdatedAt = s.clock.Now().UTC()
	if err := s.postRepository.Update(ctx, tx, post); err != nil {
		return entities.Post{}, err
	}
	return post, nil
}

func (s *PostService) Delete(
	ctx context.Context,
	tx dtx.TX,
	del entities.PostDelete,
) (entities.Post, error) {
	post, err := s.postRepository.Get(ctx, del.ID)
	if err != nil {
		return entities.Post{}, err
	}
	post.DeletedAt = pointer.Of(s.clock.Now().UTC())
	if err := s.postRepository.Update(ctx, tx, post); err != nil {
		return entities.Post{}, err
	}
	return post, nil
}
