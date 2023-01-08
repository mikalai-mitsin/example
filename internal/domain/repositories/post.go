package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/post_mock.go github.com/018bf/example/internal/domain/repositories PostRepository

type PostRepository interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Post, error)
	List(
		ctx context.Context,
		filter *models.PostFilter,
	) ([]*models.Post, error)
	Count(
		ctx context.Context,
		filter *models.PostFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		post *models.Post,
	) error
	Update(
		ctx context.Context,
		post *models.Post,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}
