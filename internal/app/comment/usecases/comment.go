package usecases

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type CommentUseCase struct {
	commentService commentService
	logger         logger
}

func NewCommentUseCase(commentService commentService, logger logger) *CommentUseCase {
	return &CommentUseCase{commentService: commentService, logger: logger}
}

func (i *CommentUseCase) Create(
	ctx context.Context,
	create entities.CommentCreate,
) (entities.Comment, error) {
	comment, err := i.commentService.Create(ctx, create)
	if err != nil {
		return entities.Comment{}, err
	}
	return comment, nil
}
func (i *CommentUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Comment, error) {
	comment, err := i.commentService.Get(ctx, id)
	if err != nil {
		return entities.Comment{}, err
	}
	return comment, nil
}

func (i *CommentUseCase) List(
	ctx context.Context,
	filter entities.CommentFilter,
) ([]entities.Comment, uint64, error) {
	items, count, err := i.commentService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (i *CommentUseCase) Update(
	ctx context.Context,
	update entities.CommentUpdate,
) (entities.Comment, error) {
	updated, err := i.commentService.Update(ctx, update)
	if err != nil {
		return entities.Comment{}, err
	}
	return updated, nil
}
func (i *CommentUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.commentService.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
