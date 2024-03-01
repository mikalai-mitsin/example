package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/day/models"
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

// DayUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/usecase.go . DayUseCase
type DayUseCase interface {
	Create(ctx context.Context, create *models.DayCreate) (*models.Day, error)
	List(ctx context.Context, filter *models.DayFilter) ([]*models.Day, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Day, error)
	Update(ctx context.Context, update *models.DayUpdate) (*models.Day, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
