package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/equipment/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// EquipmentRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . EquipmentRepository
type EquipmentRepository interface {
	Create(ctx context.Context, equipment *models.Equipment) error
	List(ctx context.Context, filter *models.EquipmentFilter) ([]*models.Equipment, error)
	Count(ctx context.Context, filter *models.EquipmentFilter) (uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Equipment, error)
	Update(ctx context.Context, equipment *models.Equipment) error
	Delete(ctx context.Context, id uuid.UUID) error
}
