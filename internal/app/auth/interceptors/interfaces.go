package interceptors

import (
	"context"
	"time"

	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	userModels "github.com/mikalai-mitsin/example/internal/app/user/models"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

// Clock - clock interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/clock.go . Clock
type Clock interface {
	Now() time.Time
}

// Logger - base logger interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/logger.go . Logger
type Logger interface {
	Debug(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
	Print(msg string, fields ...log.Field)
	Warn(msg string, fields ...log.Field)
	Error(msg string, fields ...log.Field)
	Fatal(msg string, fields ...log.Field)
	Panic(msg string, fields ...log.Field)
}

// AuthUseCase - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_usecase.go . AuthUseCase
type AuthUseCase interface {
	CreateToken(ctx context.Context, login *models.Login) (*models.TokenPair, error)
	CreateTokenByUser(ctx context.Context, user *userModels.User) (*models.TokenPair, error)
	RefreshToken(ctx context.Context, refresh models.Token) (*models.TokenPair, error)
	ValidateToken(ctx context.Context, access models.Token) error
	Auth(ctx context.Context, access models.Token) (*userModels.User, error)
}
