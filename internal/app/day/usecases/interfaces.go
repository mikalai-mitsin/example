package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/day/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// DayRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . DayRepository
type DayRepository interface {
	Create(ctx context.Context, day *models.Day) error
	List(ctx context.Context, filter *models.DayFilter) ([]*models.Day, error)
	Count(ctx context.Context, filter *models.DayFilter) (uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Day, error)
	Update(ctx context.Context, day *models.Day) error
	Delete(ctx context.Context, id uuid.UUID) error
}
