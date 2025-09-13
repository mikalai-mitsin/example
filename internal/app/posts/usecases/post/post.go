package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostUseCase struct {
	postService      postService
	postEventService postEventService
	dtxManager       dtxManager
	logger           logger
}

func NewPostUseCase(
	postService postService,
	postEventService postEventService,
	dtxManager dtxManager,
	logger logger,
) *PostUseCase {
	return &PostUseCase{
		postService:      postService,
		postEventService: postEventService,
		dtxManager:       dtxManager,
		logger:           logger,
	}
}

func (u *PostUseCase) Create(
	ctx context.Context,
	create entities.PostCreate,
) (entities.Post, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	post, err := u.postService.Create(ctx, tx, create)
	if err != nil {
		return entities.Post{}, err
	}
	if err := u.postEventService.Send(ctx, tx, post); err != nil {
		return entities.Post{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Post{}, err
	}
	return post, nil
}
func (u *PostUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Post, error) {
	post, err := u.postService.Get(ctx, id)
	if err != nil {
		return entities.Post{}, err
	}
	return post, nil
}

func (u *PostUseCase) List(
	ctx context.Context,
	filter entities.PostFilter,
) ([]entities.Post, uint64, error) {
	posts, count, err := u.postService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

func (u *PostUseCase) Update(
	ctx context.Context,
	update entities.PostUpdate,
) (entities.Post, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	post, err := u.postService.Update(ctx, tx, update)
	if err != nil {
		return entities.Post{}, err
	}
	if err := u.postEventService.Send(ctx, tx, post); err != nil {
		return entities.Post{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Post{}, err
	}
	return post, nil
}
func (u *PostUseCase) Delete(ctx context.Context, id uuid.UUID) (entities.Post, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	post, err := u.postService.Delete(ctx, tx, id)
	if err != nil {
		return entities.Post{}, err
	}
	if err := u.postEventService.Send(ctx, tx, post); err != nil {
		return entities.Post{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Post{}, err
	}
	return post, nil
}
