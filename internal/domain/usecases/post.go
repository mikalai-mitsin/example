package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/post_mock.go github.com/018bf/example/internal/domain/usecases PostUseCase

type PostUseCase interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Post, error)
	List(
		ctx context.Context,
		filter *models.PostFilter,
	) ([]*models.Post, uint64, error)
	Create(
		ctx context.Context,
		create *models.PostCreate,
	) (*models.Post, error)
	Update(
		ctx context.Context,
		update *models.PostUpdate,
	) (*models.Post, error)
	Delete(
		ctx context.Context,
		id string,
	) error
}
