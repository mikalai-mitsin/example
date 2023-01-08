package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/comment_mock.go github.com/018bf/example/internal/domain/repositories CommentRepository

type CommentRepository interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Comment, error)
	List(
		ctx context.Context,
		filter *models.CommentFilter,
	) ([]*models.Comment, error)
	Count(
		ctx context.Context,
		filter *models.CommentFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		comment *models.Comment,
	) error
	Update(
		ctx context.Context,
		comment *models.Comment,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}
