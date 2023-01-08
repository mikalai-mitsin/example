package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type PostUseCase struct {
	postRepository repositories.PostRepository
	clock          clock.Clock
	logger         log.Logger
}

func NewPostUseCase(
	postRepository repositories.PostRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.PostUseCase {
	return &PostUseCase{
		postRepository: postRepository,
		clock:          clock,
		logger:         logger,
	}
}

func (u *PostUseCase) Get(
	ctx context.Context,
	id string,
) (*models.Post, error) {
	post, err := u.postRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (u *PostUseCase) List(
	ctx context.Context,
	filter *models.PostFilter,
) ([]*models.Post, uint64, error) {
	posts, err := u.postRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.postRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

func (u *PostUseCase) Create(
	ctx context.Context,
	create *models.PostCreate,
) (*models.Post, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	post := &models.Post{
		ID:        "",
		Body:      create.Body,
		Title:     create.Title,
		UserId:    create.UserId,
		Weight:    create.Weight,
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := u.postRepository.Create(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (u *PostUseCase) Update(
	ctx context.Context,
	update *models.PostUpdate,
) (*models.Post, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	post, err := u.postRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if update.Body != nil {
		post.Body = *update.Body
	}
	if update.Title != nil {
		post.Title = *update.Title
	}
	if update.UserId != nil {
		post.UserId = *update.UserId
	}
	if update.Weight != nil {
		post.Weight = *update.Weight
	}
	post.UpdatedAt = u.clock.Now()
	if err := u.postRepository.Update(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (u *PostUseCase) Delete(ctx context.Context, id string) error {
	if err := u.postRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
