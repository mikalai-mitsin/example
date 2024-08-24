package usecases

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	userModels "github.com/mikalai-mitsin/example/internal/app/user/models"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

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
