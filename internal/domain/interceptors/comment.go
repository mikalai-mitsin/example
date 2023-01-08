package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/comment_mock.go github.com/018bf/example/internal/domain/interceptors CommentInterceptor

type CommentInterceptor interface {
	Get(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) (*models.Comment, error)
	List(
		ctx context.Context,
		filter *models.CommentFilter,
		requestUser *models.User,
	) ([]*models.Comment, uint64, error)
	Create(
		ctx context.Context,
		create *models.CommentCreate,
		requestUser *models.User,
	) (*models.Comment, error)
	Update(
		ctx context.Context,
		update *models.CommentUpdate,
		requestUser *models.User,
	) (*models.Comment, error)
	Delete(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) error
}
