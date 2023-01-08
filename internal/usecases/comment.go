package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type CommentUseCase struct {
	commentRepository repositories.CommentRepository
	clock             clock.Clock
	logger            log.Logger
}

func NewCommentUseCase(
	commentRepository repositories.CommentRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.CommentUseCase {
	return &CommentUseCase{
		commentRepository: commentRepository,
		clock:             clock,
		logger:            logger,
	}
}

func (u *CommentUseCase) Get(
	ctx context.Context,
	id string,
) (*models.Comment, error) {
	comment, err := u.commentRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (u *CommentUseCase) List(
	ctx context.Context,
	filter *models.CommentFilter,
) ([]*models.Comment, uint64, error) {
	comments, err := u.commentRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.commentRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return comments, count, nil
}

func (u *CommentUseCase) Create(
	ctx context.Context,
	create *models.CommentCreate,
) (*models.Comment, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	comment := &models.Comment{
		ID:        "",
		Body:      create.Body,
		PostId:    create.PostId,
		UserId:    create.UserId,
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := u.commentRepository.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (u *CommentUseCase) Update(
	ctx context.Context,
	update *models.CommentUpdate,
) (*models.Comment, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	comment, err := u.commentRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if update.Body != nil {
		comment.Body = *update.Body
	}
	if update.PostId != nil {
		comment.PostId = *update.PostId
	}
	if update.UserId != nil {
		comment.UserId = *update.UserId
	}
	comment.UpdatedAt = u.clock.Now()
	if err := u.commentRepository.Update(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (u *CommentUseCase) Delete(ctx context.Context, id string) error {
	if err := u.commentRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
