package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/arch/models"
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

// ArchUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/usecase.go . ArchUseCase
type ArchUseCase interface {
	Create(ctx context.Context, create *models.ArchCreate) (*models.Arch, error)
	List(ctx context.Context, filter *models.ArchFilter) ([]*models.Arch, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Arch, error)
	Update(ctx context.Context, update *models.ArchUpdate) (*models.Arch, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
