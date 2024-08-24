package usecases

import (
	"context"
	"time"

	"github.com/mikalai-mitsin/example/internal/app/user/models"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

// UserRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/repository.go . UserRepository
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	List(ctx context.Context, filter *models.UserFilter) ([]*models.User, error)
	Count(ctx context.Context, filter *models.UserFilter) (uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

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

// UUIDGenerator - UUID generator interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/uuid_generator.go . UUIDGenerator
type UUIDGenerator interface {
	NewUUID() uuid.UUID
}
