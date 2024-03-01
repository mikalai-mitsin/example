package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/day/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// DayInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . DayInterceptor
type DayInterceptor interface {
	Create(ctx context.Context, create *models.DayCreate) (*models.Day, error)
	List(ctx context.Context, filter *models.DayFilter) ([]*models.Day, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Day, error)
	Update(ctx context.Context, update *models.DayUpdate) (*models.Day, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
