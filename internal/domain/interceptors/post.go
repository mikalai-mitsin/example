package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/post_mock.go github.com/018bf/example/internal/domain/interceptors PostInterceptor

type PostInterceptor interface {
	Get(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) (*models.Post, error)
	List(
		ctx context.Context,
		filter *models.PostFilter,
		requestUser *models.User,
	) ([]*models.Post, uint64, error)
	Create(
		ctx context.Context,
		create *models.PostCreate,
		requestUser *models.User,
	) (*models.Post, error)
	Update(
		ctx context.Context,
		update *models.PostUpdate,
		requestUser *models.User,
	) (*models.Post, error)
	Delete(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) error
}
