package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/plan/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// PlanInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . PlanInterceptor
type PlanInterceptor interface {
	Create(ctx context.Context, create *models.PlanCreate) (*models.Plan, error)
	List(ctx context.Context, filter *models.PlanFilter) ([]*models.Plan, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Plan, error)
	Update(ctx context.Context, update *models.PlanUpdate) (*models.Plan, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
