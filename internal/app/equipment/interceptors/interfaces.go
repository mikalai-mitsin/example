package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/equipment/models"
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

// EquipmentUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/usecase.go . EquipmentUseCase
type EquipmentUseCase interface {
	Create(ctx context.Context, create *models.EquipmentCreate) (*models.Equipment, error)
	List(ctx context.Context, filter *models.EquipmentFilter) ([]*models.Equipment, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Equipment, error)
	Update(ctx context.Context, update *models.EquipmentUpdate) (*models.Equipment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
