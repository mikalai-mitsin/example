package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/arch/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// ArchRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . ArchRepository
type ArchRepository interface {
	Create(ctx context.Context, arch *models.Arch) error
	List(ctx context.Context, filter *models.ArchFilter) ([]*models.Arch, error)
	Count(ctx context.Context, filter *models.ArchFilter) (uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Arch, error)
	Update(ctx context.Context, arch *models.Arch) error
	Delete(ctx context.Context, id uuid.UUID) error
}
