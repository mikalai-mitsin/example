package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// UserInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . UserInterceptor
type UserInterceptor interface {
	Create(ctx context.Context, create *models.UserCreate) (*models.User, error)
	List(ctx context.Context, filter *models.UserFilter) ([]*models.User, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, update *models.UserUpdate) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
