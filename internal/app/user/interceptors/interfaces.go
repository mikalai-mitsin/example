package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/user/models"
	userModels "github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// AuthUseCase - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth.go . AuthUseCase
type AuthUseCase interface {
	GetUser(ctx context.Context) (*userModels.User, error)
	HasPermission(
		ctx context.Context,
		user *userModels.User,
		permission userModels.PermissionID,
	) error
	HasObjectPermission(
		ctx context.Context,
		user *userModels.User,
		permission userModels.PermissionID,
		object any,
	) error
}

// UserUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/usecase.go . UserUseCase
type UserUseCase interface {
	Create(ctx context.Context, create *models.UserCreate) (*models.User, error)
	List(ctx context.Context, filter *models.UserFilter) ([]*models.User, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, update *models.UserUpdate) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
