package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/auth/models"
	userModels "github.com/018bf/example/internal/app/user/models"
)

// AuthInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_interceptor.go . AuthInterceptor
type AuthInterceptor interface {
	Auth(ctx context.Context, token models.Token) (*userModels.User, error)
}
