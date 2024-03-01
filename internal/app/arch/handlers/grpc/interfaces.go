package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/arch/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// ArchInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . ArchInterceptor
type ArchInterceptor interface {
	Create(ctx context.Context, create *models.ArchCreate) (*models.Arch, error)
	List(ctx context.Context, filter *models.ArchFilter) ([]*models.Arch, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Arch, error)
	Update(ctx context.Context, update *models.ArchUpdate) (*models.Arch, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
