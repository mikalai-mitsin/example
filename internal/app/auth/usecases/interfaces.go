package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/auth/models"
	userModels "github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// AuthRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_repository.go . AuthRepository
type AuthRepository interface {
	Create(ctx context.Context, user *userModels.User) (*models.TokenPair, error)
	Validate(ctx context.Context, token models.Token) error
	RefreshToken(ctx context.Context, token models.Token) (*models.TokenPair, error)
	GetSubject(ctx context.Context, token models.Token) (string, error)
}

// UserRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_repository.go . UserRepository
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*userModels.User, error)
	Get(ctx context.Context, id uuid.UUID) (*userModels.User, error)
}

// PermissionRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/permission_repository.go . PermissionRepository
type PermissionRepository interface {
	HasPermission(
		ctx context.Context,
		permission userModels.PermissionID,
		requestUser *userModels.User,
	) error
	HasObjectPermission(
		ctx context.Context,
		permission userModels.PermissionID,
		user *userModels.User,
		obj any,
	) error
}
