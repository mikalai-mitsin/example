package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/comment_mock.go github.com/018bf/example/internal/domain/usecases CommentUseCase

type CommentUseCase interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Comment, error)
	List(
		ctx context.Context,
		filter *models.CommentFilter,
	) ([]*models.Comment, uint64, error)
	Create(
		ctx context.Context,
		create *models.CommentCreate,
	) (*models.Comment, error)
	Update(
		ctx context.Context,
		update *models.CommentUpdate,
	) (*models.Comment, error)
	Delete(
		ctx context.Context,
		id string,
	) error
}
