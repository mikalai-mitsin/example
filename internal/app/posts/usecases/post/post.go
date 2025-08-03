package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostUseCase struct {
	postService postService
	logger      logger
}

func NewPostUseCase(postService postService, logger logger) *PostUseCase {
	return &PostUseCase{postService: postService, logger: logger}
}

func (i *PostUseCase) Create(
	ctx context.Context,
	create entities.PostCreate,
) (entities.Post, error) {
	post, err := i.postService.Create(ctx, create)
	if err != nil {
		return entities.Post{}, err
	}
	return post, nil
}
func (i *PostUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Post, error) {
	post, err := i.postService.Get(ctx, id)
	if err != nil {
		return entities.Post{}, err
	}
	return post, nil
}

func (i *PostUseCase) List(
	ctx context.Context,
	filter entities.PostFilter,
) ([]entities.Post, uint64, error) {
	items, count, err := i.postService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (i *PostUseCase) Update(
	ctx context.Context,
	update entities.PostUpdate,
) (entities.Post, error) {
	updated, err := i.postService.Update(ctx, update)
	if err != nil {
		return entities.Post{}, err
	}
	return updated, nil
}
func (i *PostUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.postService.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
