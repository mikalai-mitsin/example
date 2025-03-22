package services

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/post/entities"
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

func (u *PostService) Create(
	ctx context.Context,
	create entities.PostCreate,
) (entities.Post, error) {
	if err := create.Validate(); err != nil {
		return entities.Post{}, err
	}
	now := u.clock.Now().UTC()
	post := entities.Post{
		ID:          u.uuid.NewUUID(),
		UpdatedAt:   now,
		CreatedAt:   now,
		Title:       create.Title,
		Body:        create.Body,
		IsPrivate:   create.IsPrivate,
		Tags:        create.Tags,
		PublishedAt: create.PublishedAt,
		AuthorId:    create.AuthorId,
	}
	if err := u.postRepository.Create(ctx, post); err != nil {
		return entities.Post{}, err
	}
	return post, nil
}
func (u *PostService) Get(ctx context.Context, id uuid.UUID) (entities.Post, error) {
	post, err := u.postRepository.Get(ctx, id)
	if err != nil {
		return entities.Post{}, err
	}
	return post, nil
}

func (u *PostService) List(
	ctx context.Context,
	filter entities.PostFilter,
) ([]entities.Post, uint64, error) {
	post, err := u.postRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.postRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return post, count, nil
}

func (u *PostService) Update(
	ctx context.Context,
	update entities.PostUpdate,
) (entities.Post, error) {
	if err := update.Validate(); err != nil {
		return entities.Post{}, err
	}
	post, err := u.postRepository.Get(ctx, update.ID)
	if err != nil {
		return entities.Post{}, err
	}
	{
		if update.Title != nil {
			post.Title = *update.Title
		}
		if update.Body != nil {
			post.Body = *update.Body
		}
		if update.IsPrivate != nil {
			post.IsPrivate = *update.IsPrivate
		}
		if update.Tags != nil {
			post.Tags = *update.Tags
		}
		if update.PublishedAt != nil {
			post.PublishedAt = *update.PublishedAt
		}
		if update.AuthorId != nil {
			post.AuthorId = *update.AuthorId
		}
	}
	post.UpdatedAt = u.clock.Now().UTC()
	if err := u.postRepository.Update(ctx, post); err != nil {
		return entities.Post{}, err
	}
	return post, nil
}
func (u *PostService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.postRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
