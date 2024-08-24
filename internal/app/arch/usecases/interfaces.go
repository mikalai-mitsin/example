package usecases

import (
	"context"
	"time"

	"github.com/mikalai-mitsin/example/internal/app/arch/models"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

// ArchRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/repository.go . ArchRepository
type ArchRepository interface {
	Create(ctx context.Context, arch *models.Arch) error
	List(ctx context.Context, filter *models.ArchFilter) ([]*models.Arch, error)
	Count(ctx context.Context, filter *models.ArchFilter) (uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Arch, error)
	Update(ctx context.Context, arch *models.Arch) error
	Delete(ctx context.Context, id uuid.UUID) error
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
