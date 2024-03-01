package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/plan/models"
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

// PlanUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/usecase.go . PlanUseCase
type PlanUseCase interface {
	Create(ctx context.Context, create *models.PlanCreate) (*models.Plan, error)
	List(ctx context.Context, filter *models.PlanFilter) ([]*models.Plan, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Plan, error)
	Update(ctx context.Context, update *models.PlanUpdate) (*models.Plan, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
