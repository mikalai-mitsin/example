package services

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type CommentService struct {
	commentRepository commentRepository
	clock             clock
	logger            logger
	uuid              uuidGenerator
}

func NewCommentService(
	commentRepository commentRepository,
	clock clock,
	logger logger,
	uuid uuidGenerator,
) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
		clock:             clock,
		logger:            logger,
		uuid:              uuid,
	}
}

func (u *CommentService) Create(
	ctx context.Context,
	create entities.CommentCreate,
) (entities.Comment, error) {
	if err := create.Validate(); err != nil {
		return entities.Comment{}, err
	}
	now := u.clock.Now().UTC()
	comment := entities.Comment{
		ID:        u.uuid.NewUUID(),
		UpdatedAt: now,
		CreatedAt: now,
		Text:      create.Text,
		AuthorId:  create.AuthorId,
		PostId:    create.PostId,
	}
	if err := u.commentRepository.Create(ctx, comment); err != nil {
		return entities.Comment{}, err
	}
	return comment, nil
}
func (u *CommentService) Get(ctx context.Context, id uuid.UUID) (entities.Comment, error) {
	comment, err := u.commentRepository.Get(ctx, id)
	if err != nil {
		return entities.Comment{}, err
	}
	return comment, nil
}

func (u *CommentService) List(
	ctx context.Context,
	filter entities.CommentFilter,
) ([]entities.Comment, uint64, error) {
	comment, err := u.commentRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.commentRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return comment, count, nil
}

func (u *CommentService) Update(
	ctx context.Context,
	update entities.CommentUpdate,
) (entities.Comment, error) {
	if err := update.Validate(); err != nil {
		return entities.Comment{}, err
	}
	comment, err := u.commentRepository.Get(ctx, update.ID)
	if err != nil {
		return entities.Comment{}, err
	}
	{
		if update.Text != nil {
			comment.Text = *update.Text
		}
		if update.AuthorId != nil {
			comment.AuthorId = *update.AuthorId
		}
		if update.PostId != nil {
			comment.PostId = *update.PostId
		}
	}
	comment.UpdatedAt = u.clock.Now().UTC()
	if err := u.commentRepository.Update(ctx, comment); err != nil {
		return entities.Comment{}, err
	}
	return comment, nil
}
func (u *CommentService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.commentRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
