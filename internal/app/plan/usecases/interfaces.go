package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/plan/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// PlanRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . PlanRepository
type PlanRepository interface {
	Create(ctx context.Context, plan *models.Plan) error
	List(ctx context.Context, filter *models.PlanFilter) ([]*models.Plan, error)
	Count(ctx context.Context, filter *models.PlanFilter) (uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Plan, error)
	Update(ctx context.Context, plan *models.Plan) error
	Delete(ctx context.Context, id uuid.UUID) error
}
