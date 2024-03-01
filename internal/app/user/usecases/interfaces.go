package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// UserRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . UserRepository
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	List(ctx context.Context, filter *models.UserFilter) ([]*models.User, error)
	Count(ctx context.Context, filter *models.UserFilter) (uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
