package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/equipment/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// EquipmentInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . EquipmentInterceptor
type EquipmentInterceptor interface {
	Create(ctx context.Context, create *models.EquipmentCreate) (*models.Equipment, error)
	List(ctx context.Context, filter *models.EquipmentFilter) ([]*models.Equipment, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Equipment, error)
	Update(ctx context.Context, update *models.EquipmentUpdate) (*models.Equipment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
