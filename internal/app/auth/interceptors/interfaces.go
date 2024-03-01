package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/auth/models"
	userModels "github.com/018bf/example/internal/app/user/models"
)

// AuthUseCase - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_usecase.go . AuthUseCase
type AuthUseCase interface {
	CreateToken(ctx context.Context, login *models.Login) (*models.TokenPair, error)
	CreateTokenByUser(ctx context.Context, user *userModels.User) (*models.TokenPair, error)
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
	RefreshToken(ctx context.Context, refresh models.Token) (*models.TokenPair, error)
	ValidateToken(ctx context.Context, access models.Token) error
	Auth(ctx context.Context, access models.Token) (*userModels.User, error)
}
